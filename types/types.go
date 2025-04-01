package types

type EmailPayload struct {
	Email string `json:"hashed_email"`
}

type PasswordPayload struct {
	Password string `json:"hashed_password"`
}

type Password struct {
	PasswordHash string `json:"password_hash"`
}

type Email struct {
	EmailHash int `json:"email_hash"`
}

type PasswordSuffix struct {
	Suffix string
	Count  int // número de ocorrências nos vazamentos
}

type HashPrefix struct {
	Prefix   string
	Suffixes []PasswordSuffix
}

type PasswordStore interface {
	ProcessPasswordHashes(searchPrefix string) (*HashPrefix, error)
}

type EmailStore interface {
}
