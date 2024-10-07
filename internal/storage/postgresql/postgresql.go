package postgresql

import (
	"auth/internal/storage"
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Storage struct {
	DB *pgxpool.Pool
}

func NewDBPool(ctx context.Context) (*Storage, error) {
	dbPath := os.Getenv("DB_URL")

	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	dbConn, err := pgxpool.New(ctx, dbPath)

	if err != nil {
		return nil, err
	}

	return &Storage{DB: dbConn}, nil
}


func (s *Storage) CreateUser(ctx context.Context, user storage.User) (int64, error) {
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
		return 0, err
	}

	return id, nil
}


// func (s *Storage) CreateToken(ctx context.Context, token storage.Token) error {
// 	stmt := `
// 		INSERT INTO token VALUES(@access_token, @refresh_token)
// 	`

// 	args := pgx.NamedArgs{
// 		"access_token": token.Access_token,
// 		"refresh_token": token.Refresh_token,
// 	}
	
// 	_, err := s.DB.Exec(ctx, stmt, args)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }


func (s *Storage) GetUser(ctx context.Context, id int, username string) (*storage.User, error) {
	stmt := `
		SELECT * from user where id=@id, username=@username;
	`
	
	args := pgx.NamedArgs{
		"id": id,
		"username": username,
	}

	var user storage.User
	
	err := s.DB.QueryRow(ctx, stmt, args).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// TODO: create user
// TODO: create jwt
// TODO: get tokens(access, refresh)
// TODO: get user
// TODO: get all users
// TODO: 
