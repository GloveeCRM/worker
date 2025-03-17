package postgres

import (
	"fmt"
	"glovee-worker/internal/types"
)

func (db *DB) SayHello() (*types.Hi, error) {
	const operation = "repository.postgres.SayHello"

	rows, err := db.Query("SELECT 'Hello, World!'")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	defer rows.Close()

	var greeting string

	for rows.Next() {
		if err := rows.Scan(&greeting); err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
		fmt.Println(greeting)
	}

	return &types.Hi{Message: greeting}, nil
}
