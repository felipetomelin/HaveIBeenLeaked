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
		WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method %s doesnt allowed, use POST", r.Method))
		return
	}

	var payload types.PasswordPayload
	if err := ParseJSON(r, &payload); err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("error when deserialize: %w", err))
		return
	}

	if err := validatePasswordPayload(payload); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	prefix := payload.Password
	if len(prefix) != 5 {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("hash prefix needs to be exactly 5 characters"))
		return
	}

	result, err := h.passwordStore.ProcessPasswordHashes(prefix)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("error when processing password: %w", err))
		return
	}

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
		return fmt.Errorf("hash prefix is mandatory")
	}

	for _, char := range payload.Password {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return fmt.Errorf("hash prefix needs only to contains hexadecimals characters (0-9, a-f, A-F)")
		}
	}

	return nil
}
