package user

import (
	"database/sql"
	"fmt"
	"log"
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
	query := `SELECT 
	  profile.idProfile AS profileId,
	  profile.firstName,
	  profile.lastName,
	  profile.email,
	  profile.password,
	  profile.createdAt,
	  profile.refreshToken,
	  profile.type,
	  profile.address,
	  profile.lastLogin,
	  profile.phoneNumber,
	  client.idClient,
	  client.username
	FROM profile 
	JOIN client ON profile.idProfile = client.idProfile 
	WHERE profile.email = ?`

	row := s.db.QueryRow(query, email)

	u := new(types.User)
	err := row.Scan(
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
		&u.ClientId,
		&u.Username,
	)
	if err == sql.ErrNoRows {
		return nil, nil // User not found, return nil without error
	} else if err != nil {
		return nil, err // Other error
	}

	return u, nil
}

func (s *Store) GetClientIdByUsername(username string) (string, error) {
	query := `SELECT idClient FROM client WHERE username = ?`
	row := s.db.QueryRow(query, username)
	var clientId string
	err := row.Scan(&clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("client not found for user: %s", username)
		}
		return "", fmt.Errorf("error retrieving client ID: %v", err)
	}
	return clientId, nil
}

func (s *Store) GetClientInformation(idClient string) (*types.ProfilePage, error) {
	query := `
		SELECT 
			profile.firstName, 
			profile.lastName, 
			profile.email,
			profile.address,
			profile.phoneNumber,
			client.username 
		FROM profile 
		JOIN client ON profile.idProfile = client.idProfile
		WHERE client.idClient = ?
	`

	// Get followers/following counts
	following, _ := s.CountFollowing(idClient)
	followers, _ := s.CountFollowers(idClient)

	row := s.db.QueryRow(query, idClient)

	u := new(types.UserInformation)
	err := row.Scan(
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Address,
		&u.Phone,
		&u.Username,
	)
	if err == sql.ErrNoRows {
		return nil, nil // No such client found
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving client information: %v", err)
	}

	profilePage := &types.ProfilePage{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Address:   u.Address,
		Phone:     u.Phone,
		Username:  u.Username,
		Following: following,
		Followers: followers,
	}

	return profilePage, nil
}

func (s *Store) SendRequestFriend(idFriendship string, idSender string, idReceiver string) error {
	query := `INSERT INTO friendship (idFriendship, idClient1, idClient2, status)
VALUES (?, ?, ?, 'pending');`
	_, err := s.db.Exec(query, idFriendship, idSender, idReceiver)
	if err != nil {
		return fmt.Errorf("error sending friend request: %v", err)
	}
	return nil
}

func (s *Store) AcceptRequestFriend(idFriendship string) error {
	query := `UPDATE friendship SET status = 'accepted' WHERE idFriendship = ?`
	_, err := s.db.Exec(query, idFriendship)
	if err != nil {
		return fmt.Errorf("error accepting friend request: %v", err)
	}
	return nil
}

func (s *Store) CountFollowing(idClient string) (int, error) {
	query := `SELECT COUNT(*) FROM friendship
WHERE idClient1 = ? AND status = 'accepted';`
	row := s.db.QueryRow(query, idClient, idClient)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting friendships: %v", err)
	}
	return count, nil
}

func (s *Store) CountFollowers(idClient string) (int, error) {
	query := `SELECT COUNT(*) FROM friendship
WHERE idClient2 = ? AND status = 'accepted';`
	row := s.db.QueryRow(query, idClient)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting followers: %v", err)
	}
	return count, nil
}

func (s *Store) GetFriendshipRequested(idClient string) (*[]types.Friendship, error) {
	query := `SELECT client.username,friendship.idFriendship,friendship.status,friendship.createdAt FROM friendship join client on friendship.idClient1=client.idClient WHERE status = 'pending' AND idClient2 = ?`
	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, fmt.Errorf("error retrieving friendship requests: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	friendships := []types.Friendship{}
	for rows.Next() {
		friendship := types.Friendship{}
		err = rows.Scan(
			&friendship.Username,
			&friendship.IdFriendship,
			&friendship.Status,
			&friendship.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning friendship row: %v", err)
		}
		friendships = append(friendships, friendship)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over friendship rows: %v", err)
	}
	return &friendships, nil
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

	_, err := s.db.Exec(query, idUser, user.FirstName, user.LastName, user.Email, hashedPassword, user.Address, time.Now(), time.Now(), token, user.Type, user.Phone)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (s *Store) CreateClient(idUser string, idClient string, username string) error {
	log.Println("Creating client with idUser:", idUser, "idClient:", idClient, "username:", username)
	query := `INSERT INTO client (idClient,idProfile,username)VALUES (?,?,?)`

	_, err := s.db.Exec(query, idClient, idUser, username)
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
