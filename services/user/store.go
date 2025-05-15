package user

import (
	"database/sql"
	"fmt"
	"time"

	// "github.com/wael-boudissaa/zencitiBackend/services/auth"
	"github.com/wael-boudissaa/zencitiBackend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	query := `SELECT * FROM profile WHERE email = ?`
	rows, err := s.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks

	u := new(types.User)
	for rows.Next() {
		err = rows.Scan(
			&u.Id,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.Refreshtoken,
			&u.Type,
			&u.Address,
			&u.LastLogin,
			&u.Phone,
		)
		if err != nil {
			return nil, err
		}
        return u, nil
	}
	if u.Id == "" {
		return nil, fmt.Errorf("user not found")
	}else{
        return u, nil
    }
}

func (s *Store) GetUserById(user types.User) (*types.User, error) {
	query := `SELECT * FROM profile where idProfile= ?`
	rows, err := s.db.Query(query, user.Id)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (s *Store) CreateUser(user types.RegisterUser, idUser string, token string, hashedPassword string) error {
	query := `INSERT INTO profile (idProfile, firstName, lastName, email, password, address,createdAt,lastLogin, refreshToken, type,phoneNumber)
			  VALUES (?, ?, ?, ?, ?,?,?,?,?, ?,?)`

	_, err := s.db.Exec(query, idUser, user.FirstName, user.LastName, user.Email, hashedPassword, user.Address, time.Now(), time.Now(), token, user.Type,user.Phone)

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}


	return nil
}

func (s *Store) CreateClient(idUser string,idClient string ) error {
    query := `INSERT INTO client (idClient,idProfile)VALUES (?,?)`
    _, err := s.db.Exec(query, idClient,idUser )
    if err != nil {
        return fmt.Errorf("error creating client: %v", err)
    }
    return nil
}



func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Address,

		&user.Phone,
		&user.CreatedAt,
		&user.Type,
		&user.LastLogin,
		&user.Refreshtoken,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}


