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
	Suffix string `json:"suffix"`
	Count  int    `json:"leak_count"`
}

type HashPrefix struct {
	Prefix   string           `json:"prefix"`
	Suffixes []PasswordSuffix `json:"suffixes"`
}

type PasswordStore interface {
	ProcessPasswordHashes(searchPrefix string) (*HashPrefix, error)
}

type EmailStore interface {
}
