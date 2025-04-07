package db

import "database/sql"

func EnsureTasksTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		photo_id UUID PRIMARY KEY,
		parameter TEXT,
		filter TEXT,
		status TEXT
	);`
	_, err := db.Exec(query)
	return err
}
func EnsureUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id UUID PRIMARY KEY,
		username TEXT,
		password TEXT
	);`
	_, err := db.Exec(query)
	return err
}
