package webfingo

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

// User represents a user entity from the database
type User struct {
	ID        string
	Email     string
	Username  string
	RealmID   string
	RealmName string
}

// Database interface for testing
type Database interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Close() error
}

// Postgres represents a database connection
type Postgres struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(connectionString string) (Database, error) {
	// Connect to PostgreSQL using the connection string
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Verify database connection
	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		log.Printf("Database connection error: %v", err)
		return nil, err
	}

	return &Postgres{
		db: db,
	}, nil
}

// Close closes the database connection
func (db *Postgres) Close() error {
	if db.db == nil {
		return nil
	}
	return db.db.Close()
}

// GetUserByEmail retrieves a user by their email address
func (db *Postgres) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if db.db == nil {
		return nil, errors.New("database connection not established")
	}

	query := `
		SELECT u.id, u.email, u.username, u.realm_id, r.name as realm_name
		FROM user_entity u
		INNER JOIN realms r ON u.realm_id = r.id
		WHERE u.email = $1
		LIMIT 1
	`

	var user User
	err := db.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.RealmID,
		&user.RealmName,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
