package haveibeenleaked

import (
	"SerasaLeaks/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	passwordStore types.PasswordStore
}

func NewHandler(passwordStore types.PasswordStore) *Handler {
	return &Handler{passwordStore: passwordStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/email", h.handleEmail).Methods(http.MethodPost)
	router.HandleFunc("/password", h.handlePassword).Methods(http.MethodPost)
}

func (h *Handler) handleEmail(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handlePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("método %s não permitido, use POST", r.Method))
		return
	}

	// Limitar o tamanho do corpo da requisição para evitar ataques de DoS
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	// Decodificar o payload JSON
	var payload types.PasswordPayload
	if err := ParseJSON(r, &payload); err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("erro ao analisar JSON: %w", err))
		return
	}

	// Validar o payload
	if err := validatePasswordPayload(payload); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Processar o hash da senha (primeiros 5 caracteres)
	prefix := payload.Password
	if len(prefix) != 5 {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("o prefixo do hash deve ter exatamente 5 caracteres"))
		return
	}

	// Buscar os sufixos de hash no banco de dados
	result, err := h.passwordStore.ProcessPasswordHashes(prefix)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("erro ao processar a consulta"))
		return
	}

	// Verificar se encontrou resultados
	if len(result.Suffixes) == 0 {
		WriteJSON(w, http.StatusOK, &types.HashPrefix{
			Prefix:   prefix,
			Suffixes: []types.PasswordSuffix{},
		})
		return
	}

	WriteJSON(w, http.StatusOK, result)
}

func validatePasswordPayload(payload types.PasswordPayload) error {
	if payload.Password == "" {
		return fmt.Errorf("o prefixo do hash é obrigatório")
	}

	for _, char := range payload.Password {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return fmt.Errorf("o prefixo do hash deve conter apenas caracteres hexadecimais (0-9, a-f, A-F)")
		}
	}

	return nil
}
