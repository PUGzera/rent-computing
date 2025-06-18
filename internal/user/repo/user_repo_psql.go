package user_repo

import (
	"context"
	"errors"
	"fmt"
	user "rent-computing/internal/user/data"
	"strings"

	"github.com/jackc/pgx/v5"
)

type UserRepoPSQL struct {
	table string
	conn  *pgx.Conn
}

type OptionsPSQL struct {
	Table string
	Conn  *pgx.Conn
}

func NewPSQL(options OptionsPSQL) (*UserRepoPSQL, error) {
	if options.Conn == nil {
		return nil, errors.New("no PSQL connection passed in")
	}

	table := options.Table
	if strings.EqualFold(options.Table, "") {
		table = "users"
	}

	if err := createUserTable(options.Conn, table); err != nil {
		return nil, err
	}

	return &UserRepoPSQL{
		table: table,
		conn:  options.Conn,
	}, nil
}

func createUserTable(conn *pgx.Conn, table string) error {
	query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id TEXT PRIMARY KEY,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL
		password TEXT NOT NULL
		created_at TIMESTAMP NOT NULL
    );
    `, table)
	_, err := conn.Query(context.Background(), query)
	return err
}

//CRUD Operations

// Create
func (u *UserRepoPSQL) CreateUser(ctx context.Context, user user.User) error {
	query := fmt.Sprintf("INSERT INTO %s (id, username, email, password, created_at) VALUES ($1, $2, $3, $4, $5)", u.table)
	_, err := u.conn.Query(ctx, query, user.Id, user.Username, user.Email, user.Password, user.CreatedAt)
	return err
}

// Read

func (u *UserRepoPSQL) ListUsers(ctx context.Context) ([]user.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", u.table)
	var users []user.User
	rows, err := u.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user user.User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepoPSQL) GetUser(ctx context.Context, id string) (*user.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", u.table)
	var user user.User
	err := u.conn.QueryRow(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepoPSQL) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1", u.table)
	var user user.User
	err := u.conn.QueryRow(ctx, query, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepoPSQL) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1", u.table)
	var user user.User
	err := u.conn.QueryRow(ctx, query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update
func (u *UserRepoPSQL) UpdateUser(ctx context.Context, user user.User) error {
	query := fmt.Sprintf("UPDATE %s SET id=$1, username=$2, email=$3, password=$4 WHERE id=$1", u.table)
	_, err := u.conn.Query(ctx, query, user.Id, user.Username, user.Email, user.Password)
	return err
}

// Delete
func (u *UserRepoPSQL) DeleteUser(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", u.table)
	_, err := u.conn.Query(ctx, query, id)
	return err
}

func (u *UserRepoPSQL) DeleteUserByUsername(ctx context.Context, username string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE username=$1", u.table)
	_, err := u.conn.Query(ctx, query, username)
	return err
}

func (u *UserRepoPSQL) DeleteUserByEmail(ctx context.Context, email string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE email=$1", u.table)
	_, err := u.conn.Query(ctx, query, email)
	return err
}
