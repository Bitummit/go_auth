package postgresql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Bitummit/go_auth/internal/models"
	"github.com/Bitummit/go_auth/pkg/my_errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Storage struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context) (*Storage, error) {
	dbPath := os.Getenv("DB_URL")

	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()
	
	dbConn, err := pgxpool.New(ctx, dbPath)
	// defer dbConn.Close()
	if err != nil {
		return nil, err
	}

	if err := dbConn.Ping(ctx); err != nil {
		return nil, err
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
		// if errors.Is(err, pgx.ErrNoRows)
		return 0, my_errors.ErrorUserExists
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
		return nil, fmt.Errorf("user not found %w", my_errors.ErrorNotFound)
	}
	
	return &user, nil
}
