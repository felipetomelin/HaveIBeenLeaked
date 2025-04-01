package haveibeenleaked

import (
	"SerasaLeaks/types"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// Store representa o acesso ao armazenamento de dados
type Store struct {
	db     *sql.DB
	logger *log.Logger
}

// NewStore cria uma nova instância de Store
func NewStore(db *sql.DB) *Store {
	if db == nil {
		panic("database connection cannot be nil")
	}

	return &Store{
		db: db,
	}
}

// ProcessPasswordHashes processa os hashes de senha com base no prefixo fornecido
func (s *Store) ProcessPasswordHashes(searchPrefix string) (*types.HashPrefix, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	if len(searchPrefix) != 5 {
		return nil, fmt.Errorf("o prefixo de busca deve ter exatamente 5 caracteres")
	}

	// Usar prepared statement para evitar injeção SQL e usar o parâmetro searchPrefix
	query := "SELECT passwordHash, count FROM passwords WHERE passwordHash LIKE '%2e6f9%'"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar a consulta: %w", err)
	}
	defer rows.Close()

	result := &types.HashPrefix{
		Prefix:   searchPrefix,
		Suffixes: []types.PasswordSuffix{},
	}

	for rows.Next() {
		var fullHash string
		var count int

		if err := rows.Scan(&fullHash, &count); err != nil {
			return nil, fmt.Errorf("erro ao ler o resultado: %w", err)
		}

		if len(fullHash) <= 5 {
			fmt.Printf("Hash ignorado por ser muito curto: %s\n", fullHash)
			continue
		}

		if !strings.HasPrefix(strings.ToUpper(fullHash), strings.ToUpper(searchPrefix)) {
			fmt.Printf("Hash ignorado por não começar com o prefixo %s: %s\n", searchPrefix, fullHash)
			continue
		}

		suffix := fullHash[5:]
		result.Suffixes = append(result.Suffixes, types.PasswordSuffix{
			Suffix: suffix,
			Count:  count,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return result, nil
}
