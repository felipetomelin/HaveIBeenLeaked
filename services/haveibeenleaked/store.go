package haveibeenleaked

import (
	"SerasaLeaks/types"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	db     *sql.DB
	logger *log.Logger
}

func NewStore(db *sql.DB) *Store {
	if db == nil {
		panic("database connection cannot be nil")
	}

	return &Store{
		db: db,
	}
}

func (s *Store) ProcessPasswordHashes(searchPrefix string) (*types.HashPrefix, error) {
	query := "SELECT passwordHash, count(*) FROM passwords WHERE passwordHash LIKE ? GROUP BY passwordHash;"
	rows, err := s.db.Query(query, "%"+searchPrefix+"%")
	if err != nil {
		return nil, fmt.Errorf("error when execute query: %w", err)
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
			return nil, fmt.Errorf("error to deserialize after execute query: %w", err)
		}

		suffix := fullHash[5:]
		result.Suffixes = append(result.Suffixes, types.PasswordSuffix{
			Suffix: suffix,
			Count:  count,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterate over the results: %w", err)
	}

	return result, nil
}
