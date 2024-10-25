package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Bitummit/go_auth/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Storage struct {
	DB *pgxpool.Pool
}

var ErrorNotFound = errors.New("not found")
var ErrorUserExists = errors.New("user exists")


func New(ctx context.Context) (*Storage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	dbPath := os.Getenv("DB_URL")
	dbConn, err := pgxpool.New(ctx, dbPath)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	if err := dbConn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Storage{DB: dbConn}, nil
}

func (s *Storage) CreateUser(ctx context.Context, user models.User) (int64, error) {
	stmt := `
		INSERT INTO my_user(username, pass) VALUES(@username, @password) RETURNING id;
	`
	args := pgx.NamedArgs{
		"username": user.Username,
		"password": user.Password,
	}

	var id int64
	err := s.DB.QueryRow(ctx, stmt, args).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrTooManyRows) {
			return 0, fmt.Errorf("%w", ErrorUserExists)
		}
		return 0, fmt.Errorf("query inserting user %w", err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, username string) (*models.User, error) {
	stmt := `
		SELECT * from my_user where username=@username;
	`
	args := pgx.NamedArgs{
		"username": username,
	}

	var user models.User
	err := s.DB.QueryRow(ctx, stmt, args).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w", ErrorNotFound)
		}
		return nil, fmt.Errorf("query getting user%w", err)
	}
	
	return &user, nil
}
