package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Configuração da conexão com o banco de dados
	dbUser := "root"
	dbPass := "password"
	dbHost := "localhost:3306"
	dbName := "haveibeenleaked"

	// String de conexão para MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)

	// Abrir conexão com o banco de dados
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com o banco de dados: %v", err)
	}
	defer db.Close()

	// Configurar parâmetros da conexão
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar se a conexão está funcionando
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao verificar conexão com o banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso")

	// Criar e executar o servidor
	server := NewApiServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
		os.Exit(1)
	}
}
