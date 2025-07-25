package user

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	// "github.com/wael-boudissaa/zencitiBackend/services/auth"
	"github.com/wael-boudissaa/zencitiBackend/types"
	"github.com/wael-boudissaa/zencitiBackend/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) VerifyAdminRestaurantAssignment(idAdminRestaurant string) (bool, string, error) {
	query := `SELECT idRestaurant FROM restaurant WHERE idAdminRestaurant = ?`
	var idRestaurant string
	err := s.db.QueryRow(query, idAdminRestaurant).Scan(&idRestaurant)
	if err == sql.ErrNoRows {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}
	return true, idRestaurant, nil
}

func (s *Store) SetAdminLocation(idAdmin string, latitude, longitude float64) error {
	var count int
	checkQuery := `SELECT COUNT(*) FROM admin WHERE idAdmin = ?`
	err := s.db.QueryRow(checkQuery, idAdmin).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking admin existence: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("admin with ID %s not found", idAdmin)
	}

	// Update admin location
	updateQuery := `UPDATE admin SET latitude = ?, longitude = ? WHERE idAdmin = ?`
	_, err = s.db.Exec(updateQuery, latitude, longitude, idAdmin)
	if err != nil {
		return fmt.Errorf("error updating admin location: %v", err)
	}

	return nil
}

func (s *Store) GetAdminLocation(idAdmin string) (*types.AdminLocation, error) {
	query := `
        SELECT 
            latitude,
            longitude
        FROM admin 
        WHERE idAdmin = ?
    `

	row := s.db.QueryRow(query, idAdmin)
	var adminLocation types.AdminLocation

	err := row.Scan(
		&adminLocation.Latitude,
		&adminLocation.Longitude,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin with ID %s not found", idAdmin)
		}
		return nil, fmt.Errorf("error retrieving admin location: %v", err)
	}

	if adminLocation.Latitude == nil || adminLocation.Longitude == nil {
		return &types.AdminLocation{
			Latitude:    nil,
			Longitude:   nil,
			HasLocation: false,
		}, nil
	}

	adminLocation.HasLocation = true
	return &adminLocation, nil
}

func (s *Store) CreateRestaurantWithAdmin(restaurantData types.RestaurantCreation, profileData types.RegisterAdmin) (string, string, error) {
	// Start a transaction to ensure atomicity
	tx, err := s.db.Begin()
	if err != nil {
		return "", "", fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generate IDs
	idProfile, err := utils.CreateAnId()
	if err != nil {
		return "", "", fmt.Errorf("error generating profile ID: %v", err)
	}

	idAdminRestaurant, err := utils.CreateAnId()
	if err != nil {
		return "", "", fmt.Errorf("error generating admin restaurant ID: %v", err)
	}

	idRestaurant, err := utils.CreateAnId()
	if err != nil {
		return "", "", fmt.Errorf("error generating restaurant ID: %v", err)
	}

	// Hash password
	hashedPassword, err := utils.HashedPassword(profileData.Password)
	if err != nil {
		return "", "", fmt.Errorf("error hashing password: %v", err)
	}

	// Create refresh token
	token, err := utils.CreateRefreshToken(idProfile, profileData.Type)
	if err != nil {
		return "", "", fmt.Errorf("error creating refresh token: %v", err)
	}

	// Check if email already exists
	var emailCount int
	checkEmailQuery := `SELECT COUNT(*) FROM profile WHERE email = ?`
	err = tx.QueryRow(checkEmailQuery, profileData.Email).Scan(&emailCount)
	if err != nil {
		return "", "", fmt.Errorf("error checking email existence: %v", err)
	}
	if emailCount > 0 {
		return "", "", fmt.Errorf("email %s already exists", profileData.Email)
	}

	// 1. Create profile
	profileQuery := `INSERT INTO profile (idProfile, firstName, lastName, email, password, address, createdAt, lastLogin, refreshToken, type, phoneNumber)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(profileQuery,
		idProfile,
		profileData.FirstName,
		profileData.LastName,
		profileData.Email,
		string(hashedPassword),
		profileData.Address,
		time.Now(),
		time.Now(),
		token,
		profileData.Type,
		profileData.Phone,
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating profile: %v", err)
	}

	// 2. Create adminRestaurant
	adminQuery := `INSERT INTO adminRestaurant (idAdminRestaurant, idProfile) VALUES (?, ?)`
	_, err = tx.Exec(adminQuery, idAdminRestaurant, idProfile)
	if err != nil {
		return "", "", fmt.Errorf("error creating admin restaurant: %v", err)
	}

	// 3. Create restaurant
	restaurantQuery := `INSERT INTO restaurant (idRestaurant, idAdminRestaurant, name, image, longitude, latitude, description, capacity, location) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(restaurantQuery,
		idRestaurant,
		idAdminRestaurant,
		restaurantData.Name,
		restaurantData.Image,
		restaurantData.Longitude,
		restaurantData.Latitude,
		restaurantData.Description,
		restaurantData.Capacity,
		restaurantData.Location,
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating restaurant: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return "", "", fmt.Errorf("error committing transaction: %v", err)
	}

	return idRestaurant, token, nil
}

func (s *Store) IsClientAdminActivity(idProfile string) (bool, string, error) {
	// Check if user is actively assigned to an activity
	// This means they exist in adminActivity table AND are assigned to an activity
	query := `
		SELECT aa.idAdminActivity, a.idActivity 
		FROM adminActivity aa
		JOIN activity a ON aa.idAdminActivity = a.idAdminActivity
		WHERE aa.idProfile = ?
	`
	var idAdminActivity, idActivity string
	err := s.db.QueryRow(query, idProfile).Scan(&idAdminActivity, &idActivity)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil // Not an active admin (either not in adminActivity table or not assigned to any activity)
		}
		return false, "", fmt.Errorf("error checking admin activity status: %v", err)
	}
	return true, idAdminActivity, nil
}

func (s *Store) GetAdminByEmail(email string) (*types.UserAdmin, error) {
	query := `SELECT 
	  profile.idProfile AS profileId,
	  profile.firstName,
	  profile.lastName,
	  profile.email,
	  profile.createdAt,
	  profile.type,
	  profile.address,
      profile.password,
	  profile.lastLogin,
	  profile.phoneNumber,
      adminRestaurant.idAdminRestaurant,
      restaurant.idRestaurant
	FROM profile 
    join adminRestaurant ON profile.idProfile = adminRestaurant.idProfile
    join restaurant ON adminRestaurant.idAdminRestaurant = restaurant.idAdminRestaurant
	WHERE profile.email = ?`

	row := s.db.QueryRow(query, email)

	u := new(types.UserAdmin)
	err := row.Scan(
		&u.Id,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.CreatedAt,
		&u.Type,
		&u.Address,
		&u.Password,
		&u.LastLogin,
		&u.Phone,
		&u.IdAdminRestaurant,
		&u.IdRestaurant,
	)
	if err == sql.ErrNoRows {
		return nil, nil // User not found, return nil without error
	} else if err != nil {
		return nil, err // Other error
	}

	return u, nil
}

func (s *Store) UpdateClientLocation(idClient string, longitude, latitude float64) error {
	query := `UPDATE client SET longitude = ?, latitude = ? WHERE idClient = ?`
	_, err := s.db.Exec(query, longitude, latitude, idClient)
	return err
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
    query := `SELECT 
      profile.idProfile,
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
        return nil, err // Other database error
    }

    return u, nil
}

func (s *Store) GetAllClients() ([]types.ClientInfo, error) {
	query := `
        SELECT c.idClient, p.firstName, p.lastName, p.email, c.username,
               CASE WHEN aa.idAdminActivity IS NOT NULL THEN true ELSE false END as isAdminActivity
        FROM client c
        JOIN profile p ON c.idProfile = p.idProfile
        LEFT JOIN adminActivity aa ON p.idProfile = aa.idProfile
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []types.ClientInfo
	for rows.Next() {
		var client types.ClientInfo
		err := rows.Scan(&client.IdClient, &client.FirstName, &client.LastName,
			&client.Email, &client.Username, &client.IsAdminActivity)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func (s *Store) AssignClientToAdminActivity(idClient string) error {
	var idProfile string
	err := s.db.QueryRow(`SELECT idProfile FROM client WHERE idClient = ?`, idClient).Scan(&idProfile)
	if err != nil {
		return fmt.Errorf("error getting client profile: %v", err)
	}

	var count int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM adminActivity WHERE idProfile = ?`, idProfile).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking existing admin: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("client is already an admin activity")
	}

	// Create new adminActivity
	idAdminActivity, err := utils.CreateAnId()
	if err != nil {
		return err
	}

	query := `INSERT INTO adminActivity (idAdminActivity, idProfile) VALUES (?, ?)`
	_, err = s.db.Exec(query, idAdminActivity, idProfile)
	if err != nil {
		return fmt.Errorf("error creating admin activity: %v", err)
	}
	query = `UPDATE profile set type = 'adminActivity' WHERE idProfile = ?`
	_, err = s.db.Exec(query, idProfile)
	if err != nil {
		return fmt.Errorf("error updating profile type to adminActivity: %v", err)
	}

	return nil
}

func (s *Store) SearchUsersByUsernamePrefix(prefix string) (*[]string, error) {
	query := `SELECT username FROM client WHERE username LIKE ? LIMIT 5`
	likePattern := prefix + "%"

	rows, err := s.db.Query(query, likePattern)
	if err != nil {
		return nil, fmt.Errorf("error searching usernames: %v", err)
	}
	defer rows.Close()

	usernames := []string{}
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("error scanning username: %v", err)
		}
		usernames = append(usernames, username)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating username rows: %v", err)
	}
	return &usernames, nil
}

func (s *Store) GetFriendsOfClient(idClient string) (*[]string, error) {
	query := `select idClient2 from friendship where idClient1=? and status ="accepted"`
	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, fmt.Errorf("error retrieving friends: %v", err)
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	friends := []string{}
	for rows.Next() {
		var friendId string
		err = rows.Scan(&friendId)
		if err != nil {
			return nil, fmt.Errorf("error scanning friend row: %v", err)
		}
		friends = append(friends, friendId)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over friend rows: %v", err)
	}
	return &friends, nil
}

// func convertToInterfaceSlice(strs []string) []interface{} {
// 	result := make([]interface{}, len(strs))
// 	for i, s := range strs {
// 		result[i] = s
// 	}
// 	return result
// }
//

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

func (s *Store) GetClientInformationUsername(username string) (*types.ProfilePage, error) {
	query := `
		SELECT 
        client.idClient,
			profile.firstName, 
			profile.lastName, 
			profile.email,
			profile.address,
			profile.phoneNumber,
			client.username ,
            client.following,
            client.followers
		FROM profile 
		JOIN client ON profile.idProfile = client.idProfile
		WHERE client.username = ?
	`

	row := s.db.QueryRow(query, username)
	u := new(types.UserInformation)
	err := row.Scan(
		&u.IdClient,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Address,
		&u.Phone,
		&u.Username,
		&u.Following,
		&u.Followers,
	)
	if err == sql.ErrNoRows {
		return nil, nil // No such client found
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving client information: %v", err)
	}
	// Get followers/following counts
	following, _ := s.CountFollowing(u.IdClient)
	followers, _ := s.CountFollowers(u.IdClient)

	err = s.UpdateFollowingFollowers(u.IdClient, following, followers)
	if err != nil {
		return nil, fmt.Errorf("error updating following/followers count: %v", err)
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

func (s *Store) UpdateFollowingFollowers(idClient string, following int, followers int) error {
	query := `UPDATE client SET following = ?, followers = ? WHERE idClient = ?`
	_, err := s.db.Exec(query, following, followers, idClient)
	if err != nil {
		return fmt.Errorf("error updating following/followers count: %v", err)
	}
	return nil
}

func (s *Store) GetClientInformation(idClient string) (*types.ProfilePage, error) {
	query := `
		SELECT 
        client.idClient,
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

	row := s.db.QueryRow(query, idClient)
	u := new(types.UserInformation)
	err := row.Scan(
		&u.IdClient,
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
	// Get followers/following counts
	following, _ := s.CountFollowing(u.IdClient)
	followers, _ := s.CountFollowers(u.IdClient)
	log.Println("Following count:", following, "Followers count:", followers)

	err = s.UpdateFollowingFollowers(u.IdClient, following, followers)
	if err != nil {
		return nil, fmt.Errorf("error updating following/followers count: %v", err)
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
	row := s.db.QueryRow(query, idClient)
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

func (s *Store) CreateUser(user interface{}, idUser string, token string, hashedPassword string) error {
	switch u := user.(type) {
	case types.RegisterUser:
		query := `INSERT INTO profile (idProfile, firstName, lastName, email, password, address, createdAt, lastLogin, refreshToken, type, phoneNumber)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err := s.db.Exec(query, idUser, u.FirstName, u.LastName, u.Email, hashedPassword, u.Address, time.Now(), time.Now(), token, u.Type, u.Phone)
		if err != nil {
			return fmt.Errorf("error creating user: %v", err)
		}
		return nil

	case types.RegisterAdmin:
		query := `INSERT INTO profile (idProfile, firstName, lastName, email, password, address, createdAt, lastLogin, refreshToken, type, phoneNumber)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err := s.db.Exec(query, idUser, u.FirstName, u.LastName, u.Email, hashedPassword, u.Address, time.Now(), time.Now(), token, u.Type, u.Phone)
		if err != nil {
			return fmt.Errorf("error creating admin: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("unsupported user type: %T", u)
	}
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

func (s *Store) CreateAdminRestaurant(idUser string, idAdminRestaurant string) error {
	query := `INSERT INTO adminRestaurant (idAdminRestaurant, idProfile) VALUES (?, ?)`
	_, err := s.db.Exec(query, idAdminRestaurant, idUser)
	if err != nil {
		return fmt.Errorf("error creating restaurant admin: %v", err)
	}
	return nil
}

func (s *Store) CreateAdminActivity(idUser string, idAdminActivity string) error {
	query := `INSERT INTO adminActivity (idAdminActivity, idProfile) VALUES (?, ?)`
	_, err := s.db.Exec(query, idAdminActivity, idUser)
	if err != nil {
		return fmt.Errorf("error creating activity admin: %v", err)
	}
	return nil
}

func (s *Store) UpdateRestaurantAdmin(idRestaurant string, idAdminRestaurant string) error {
	query := `UPDATE restaurant SET idAdminRestaurant = ? WHERE idRestaurant = ?`
	_, err := s.db.Exec(query, idAdminRestaurant, idRestaurant)
	if err != nil {
		return fmt.Errorf("error updating restaurant admin: %v", err)
	}
	return nil
}

// Add to services/user/store.go

// GetUserStats retrieves comprehensive user statistics
func (s *Store) GetUserStats() (*types.UserStats, error) {
	stats := &types.UserStats{}

	// Get total users
	totalQuery := "SELECT COUNT(*) FROM profile"
	err := s.db.QueryRow(totalQuery).Scan(&stats.TotalUsers)
	if err != nil {
		return nil, fmt.Errorf("error counting total users: %v", err)
	}

	// Get active users today (users who logged in today)
	activeQuery := "SELECT COUNT(*) FROM profile WHERE DATE(lastLogin) = CURDATE()"
	err = s.db.QueryRow(activeQuery).Scan(&stats.ActiveUsersToday)
	if err != nil {
		return nil, fmt.Errorf("error counting active users today: %v", err)
	}

	// Get new users this month
	newQuery := "SELECT COUNT(*) FROM profile WHERE MONTH(createdAt) = MONTH(CURDATE()) AND YEAR(createdAt) = YEAR(CURDATE())"
	err = s.db.QueryRow(newQuery).Scan(&stats.NewUsersThisMonth)
	if err != nil {
		return nil, fmt.Errorf("error counting new users this month: %v", err)
	}

	// Get monthly stats for the last year
	monthlyStats, err := s.GetMonthlyUserStats()
	if err != nil {
		return nil, err
	}
	stats.MonthlyStats = monthlyStats

	return stats, nil
}

// GetMonthlyUserStats retrieves monthly user statistics for the past year
func (s *Store) GetMonthlyUserStats() ([]types.MonthlyUserStats, error) {
	// Create a map to hold our results
	statMap := make(map[string]types.MonthlyUserStats)

	// First, get new users per month for the last year
	newUsersQuery := `
        SELECT 
            MONTH(createdAt) as month,
            YEAR(createdAt) as year,
            COUNT(*) as newUsers
        FROM 
            profile
        WHERE 
            createdAt >= DATE_SUB(CURDATE(), INTERVAL 1 YEAR)
        GROUP BY 
            YEAR(createdAt), MONTH(createdAt)
        ORDER BY 
            YEAR(createdAt), MONTH(createdAt)
    `

	rows, err := s.db.Query(newUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("error retrieving monthly new users: %v", err)
	}
	defer rows.Close()

	// Process new users data
	for rows.Next() {
		var month, year, newUsers int
		if err := rows.Scan(&month, &year, &newUsers); err != nil {
			return nil, fmt.Errorf("error scanning monthly new users row: %v", err)
		}
		key := fmt.Sprintf("%d-%d", year, month)
		statMap[key] = types.MonthlyUserStats{
			Month:       month,
			Year:        year,
			NewUsers:    newUsers,
			ActiveUsers: 0, // Will be updated in next query
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating monthly new users rows: %v", err)
	}

	// Next, get active users per month for the last year
	activeUsersQuery := `
        SELECT 
            MONTH(lastLogin) as month,
            YEAR(lastLogin) as year,
            COUNT(DISTINCT idProfile) as activeUsers
        FROM 
            profile
        WHERE 
            lastLogin >= DATE_SUB(CURDATE(), INTERVAL 1 YEAR)
        GROUP BY 
            YEAR(lastLogin), MONTH(lastLogin)
        ORDER BY 
            YEAR(lastLogin), MONTH(lastLogin)
    `

	activeRows, err := s.db.Query(activeUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("error retrieving monthly active users: %v", err)
	}
	defer activeRows.Close()

	// Process active users data
	for activeRows.Next() {
		var month, year, activeUsers int
		if err := activeRows.Scan(&month, &year, &activeUsers); err != nil {
			return nil, fmt.Errorf("error scanning monthly active users row: %v", err)
		}
		key := fmt.Sprintf("%d-%d", year, month)
		if stat, exists := statMap[key]; exists {
			// Update existing entry
			stat.ActiveUsers = activeUsers
			statMap[key] = stat
		} else {
			// Create new entry if month-year doesn't exist yet
			statMap[key] = types.MonthlyUserStats{
				Month:       month,
				Year:        year,
				NewUsers:    0,
				ActiveUsers: activeUsers,
			}
		}
	}
	if err := activeRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating monthly active users rows: %v", err)
	}

	// Convert map to sorted slice
	var result []types.MonthlyUserStats
	for _, stat := range statMap {
		result = append(result, stat)
	}

	// Sort by year and month
	sort.Slice(result, func(i, j int) bool {
		if result[i].Year != result[j].Year {
			return result[i].Year < result[j].Year
		}
		return result[i].Month < result[j].Month
	})

	return result, nil
}

func (s *Store) CreateActivityWithAdmin(activityData types.ActivityCreationWithAdmin, profileData types.ActivityAdminCreation) (string, string, error) {
    // Start a transaction to ensure atomicity
    tx, err := s.db.Begin()
    if err != nil {
        return "", "", fmt.Errorf("error starting transaction: %v", err)
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // Generate IDs
    idProfile, err := utils.CreateAnId()
    if err != nil {
        return "", "", fmt.Errorf("error generating profile ID: %v", err)
    }

    idAdminActivity, err := utils.CreateAnId()
    if err != nil {
        return "", "", fmt.Errorf("error generating admin activity ID: %v", err)
    }

    idActivity, err := utils.CreateAnId()
    if err != nil {
        return "", "", fmt.Errorf("error generating activity ID: %v", err)
    }

    // Generate client ID for dual role
    idClient, err := utils.CreateAnId()
    if err != nil {
        return "", "", fmt.Errorf("error generating client ID: %v", err)
    }

    // Generate username (you might want to make this configurable)
    username := fmt.Sprintf("admin_%s", idAdminActivity[:8])

    // Hash password
    hashedPassword, err := utils.HashedPassword(profileData.Password)
    if err != nil {
        return "", "", fmt.Errorf("error hashing password: %v", err)
    }

    // Create refresh token
    token, err := utils.CreateRefreshToken(idProfile, profileData.Type)
    if err != nil {
        return "", "", fmt.Errorf("error creating refresh token: %v", err)
    }

    // Check if email already exists
    var emailCount int
    checkEmailQuery := `SELECT COUNT(*) FROM profile WHERE email = ?`
    err = tx.QueryRow(checkEmailQuery, profileData.Email).Scan(&emailCount)
    if err != nil {
        return "", "", fmt.Errorf("error checking email existence: %v", err)
    }
    if emailCount > 0 {
        return "", "", fmt.Errorf("email %s already exists", profileData.Email)
    }

    // 1. Create profile
    profileQuery := `INSERT INTO profile (idProfile, firstName, lastName, email, password, address, createdAt, lastLogin, refreshToken, type, phoneNumber)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

    _, err = tx.Exec(profileQuery,
        idProfile,
        profileData.FirstName,
        profileData.LastName,
        profileData.Email,
        string(hashedPassword),
        profileData.Address,
        time.Now(),
        time.Now(),
        token,
        profileData.Type,
        profileData.Phone,
    )
    if err != nil {
        return "", "", fmt.Errorf("error creating profile: %v", err)
    }

    // 2. Create adminActivity
    adminQuery := `INSERT INTO adminActivity (idAdminActivity, idProfile) VALUES (?, ?)`
    _, err = tx.Exec(adminQuery, idAdminActivity, idProfile)
    if err != nil {
        return "", "", fmt.Errorf("error creating admin activity: %v", err)
    }

    // 3. Create client (so they can also be a regular client)
    clientQuery := `INSERT INTO client (idClient, idProfile, username, longitude, latitude, following, followers) VALUES (?, ?, ?, 0, 0, 0, 0)`
    _, err = tx.Exec(clientQuery, idClient, idProfile, username)
    if err != nil {
        return "", "", fmt.Errorf("error creating client: %v", err)
    }

    // 4. Create activity
    activityQuery := `INSERT INTO activity (idActivity, idAdminActivity, nameActivity, descriptionActivity, imageActivity, longitude, latitude, idTypeActivity, capacity) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
    _, err = tx.Exec(activityQuery,
        idActivity,
        idAdminActivity,
        activityData.Name,
        activityData.Description,
        activityData.Image,
        activityData.Longitude,
        activityData.Latitude,
        activityData.IdTypeActivity,
        activityData.Capacity, // Use capacity from form data
    )
    if err != nil {
        return "", "", fmt.Errorf("error creating activity: %v", err)
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        return "", "", fmt.Errorf("error committing transaction: %v", err)
    }

    return idActivity, token, nil
}

func (s *Store) GetGeneralAdminByEmail(email string) (*types.User, error) {
    query := `
        SELECT p.idProfile, p.firstName, p.lastName, p.email, p.password, p.type, p.address, p.phoneNumber
        FROM profile p
        WHERE p.email = ? AND p.type = 'admin'
    `
    
    var u types.User
    err := s.db.QueryRow(query, email).Scan(
        &u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Type, &u.Address, &u.Phone,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("error getting admin by email: %v", err)
    }
    
    return &u, nil
}

func (s *Store) GetAllCampusUsers() ([]types.CampusUser, error) {
    query := `
        SELECT 
            p.idProfile, p.firstName, p.lastName, p.email, p.type, p.address, p.phoneNumber, p.createdAt,
            c.idClient, c.username,
            a.idAdmin,
            aa.idAdminActivity,
            ar.idAdminRestaurant,
            -- Check if adminActivity is actively assigned to an activity
            act.idActivity as assignedActivityId,
            act.nameActivity as assignedActivityName,
            -- Check if adminRestaurant is actively assigned to a restaurant
            rest.idRestaurant as assignedRestaurantId,
            rest.name as assignedRestaurantName
        FROM profile p
        LEFT JOIN client c ON p.idProfile = c.idProfile
        LEFT JOIN admin a ON p.idProfile = a.idProfile
        LEFT JOIN adminActivity aa ON p.idProfile = aa.idProfile
        LEFT JOIN adminRestaurant ar ON p.idProfile = ar.idProfile
        -- Check if the adminActivity is currently managing an activity
        LEFT JOIN activity act ON aa.idAdminActivity = act.idAdminActivity
        -- Check if the adminRestaurant is currently managing a restaurant
        LEFT JOIN restaurant rest ON ar.idAdminRestaurant = rest.idAdminRestaurant
        ORDER BY p.createdAt DESC
    `
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error getting campus users: %v", err)
    }
    defer rows.Close()
    
    var users []types.CampusUser
    
    for rows.Next() {
        var user types.CampusUser
        var idClient, username, idAdmin, idAdminActivity, idAdminRestaurant sql.NullString
        var address, phoneNumber sql.NullString
        var assignedActivityId, assignedRestaurantId sql.NullString
        var assignedActivityName, assignedRestaurantName sql.NullString
        
        err := rows.Scan(
            &user.IdProfile, &user.FirstName, &user.LastName, &user.Email, &user.Type,
            &address, &phoneNumber, &user.CreatedAt,
            &idClient, &username, &idAdmin, &idAdminActivity, &idAdminRestaurant,
            &assignedActivityId, &assignedActivityName,
            &assignedRestaurantId, &assignedRestaurantName,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning user row: %v", err)
        }
        
        // Set nullable fields
        if address.Valid {
            user.Address = &address.String
        }
        if phoneNumber.Valid {
            user.PhoneNumber = &phoneNumber.String
        }
        if username.Valid {
            user.Username = &username.String
        }
        if idClient.Valid {
            user.IdClient = &idClient.String
        }
        if idAdmin.Valid {
            user.IdAdmin = &idAdmin.String
        }
        if idAdminActivity.Valid {
            user.IdAdminActivity = &idAdminActivity.String
        }
        if idAdminRestaurant.Valid {
            user.IdAdminRestaurant = &idAdminRestaurant.String
        }
        
        // Determine admin status and assigned entities
        if idAdminActivity.Valid {
            if assignedActivityId.Valid {
                user.AdminActivityStatus = "active"
                user.AssignedActivityId = &assignedActivityId.String
                if assignedActivityName.Valid {
                    user.AssignedActivityName = &assignedActivityName.String
                }
            } else {
                user.AdminActivityStatus = "inactive"
            }
        }
        
        if idAdminRestaurant.Valid {
            if assignedRestaurantId.Valid {
                user.AdminRestaurantStatus = "active"
                user.AssignedRestaurantId = &assignedRestaurantId.String
                if assignedRestaurantName.Valid {
                    user.AssignedRestaurantName = &assignedRestaurantName.String
                }
            } else {
                user.AdminRestaurantStatus = "inactive"
            }
        }
        
        // Determine roles - only include active admin roles
        var roles []string
        
        // Only include primary type if it's not an admin type, or if it's an admin type and the user is active
        if user.Type == "client" || user.Type == "admin" {
            roles = append(roles, user.Type)
        } else if user.Type == "adminActivity" && user.AdminActivityStatus == "active" {
            roles = append(roles, "adminActivity")
        } else if user.Type == "adminRestaurant" && user.AdminRestaurantStatus == "active" {
            roles = append(roles, "adminRestaurant")
        }
        
        // Add additional roles only if they are different from primary type and active
        if idClient.Valid && user.Type != "client" {
            roles = append(roles, "client")
        }
        if idAdminActivity.Valid && user.AdminActivityStatus == "active" && user.Type != "adminActivity" {
            roles = append(roles, "adminActivity")
        }
        if idAdminRestaurant.Valid && user.AdminRestaurantStatus == "active" && user.Type != "adminRestaurant" {
            roles = append(roles, "adminRestaurant")
        }
        if idAdmin.Valid && user.Type != "admin" {
            roles = append(roles, "admin")
        }
        
        user.Roles = roles
        users = append(users, user)
    }
    
    return users, nil
}

func (s *Store) AssignUserToRole(idUser string, role string) error {
    // First check if user exists
    var userExists bool
    checkQuery := `SELECT EXISTS(SELECT 1 FROM profile WHERE idProfile = ?)`
    err := s.db.QueryRow(checkQuery, idUser).Scan(&userExists)
    if err != nil {
        return fmt.Errorf("error checking user existence: %v", err)
    }
    if !userExists {
        return fmt.Errorf("user not found")
    }
    
    // Generate new ID for the role
    newRoleId, err := utils.CreateAnId()
    if err != nil {
        return fmt.Errorf("error generating role ID: %v", err)
    }
    
    // Assign the role based on type
    switch role {
    case "adminActivity":
        // Check if user is already ACTIVELY assigned to any activity as admin
        // This means they exist in adminActivity table AND are assigned to an activity
        var existingActiveActivityCount int
        checkActiveQuery := `
            SELECT COUNT(*) 
            FROM adminActivity aa
            JOIN activity a ON aa.idAdminActivity = a.idAdminActivity
            WHERE aa.idProfile = ?
        `
        err = s.db.QueryRow(checkActiveQuery, idUser).Scan(&existingActiveActivityCount)
        if err != nil {
            return fmt.Errorf("error checking existing active adminActivity assignment: %v", err)
        }
        if existingActiveActivityCount > 0 {
            return fmt.Errorf("user is already actively assigned as admin to an activity. An admin can only manage one activity")
        }
        
        // Check if user already has an adminActivity record, if not create one
        var existingAdminActivityId string
        checkAdminQuery := `SELECT idAdminActivity FROM adminActivity WHERE idProfile = ?`
        err = s.db.QueryRow(checkAdminQuery, idUser).Scan(&existingAdminActivityId)
        
        if err == sql.ErrNoRows {
            // User doesn't have adminActivity record, create one
            insertQuery := `INSERT INTO adminActivity (idAdminActivity, idProfile) VALUES (?, ?)`
            _, err = s.db.Exec(insertQuery, newRoleId, idUser)
            if err != nil {
                return fmt.Errorf("error creating adminActivity role: %v", err)
            }
        } else if err != nil {
            return fmt.Errorf("error checking existing adminActivity: %v", err)
        }
        
        // Note: This method doesn't assign to specific activity, just creates the admin role
        // Use AssignUserToRoleWithEntity to assign to specific activity
        
    case "adminRestaurant":
        // Check if user is already ACTIVELY assigned to any restaurant as admin
        // This means they exist in adminRestaurant table AND are assigned to a restaurant
        var existingActiveRestaurantCount int
        checkActiveQuery := `
            SELECT COUNT(*) 
            FROM adminRestaurant ar
            JOIN restaurant r ON ar.idAdminRestaurant = r.idAdminRestaurant
            WHERE ar.idProfile = ?
        `
        err = s.db.QueryRow(checkActiveQuery, idUser).Scan(&existingActiveRestaurantCount)
        if err != nil {
            return fmt.Errorf("error checking existing active adminRestaurant assignment: %v", err)
        }
        if existingActiveRestaurantCount > 0 {
            return fmt.Errorf("user is already actively assigned as admin to a restaurant. An admin can only manage one restaurant")
        }
        
        // Check if user already has an adminRestaurant record, if not create one
        var existingAdminRestaurantId string
        checkAdminQuery := `SELECT idAdminRestaurant FROM adminRestaurant WHERE idProfile = ?`
        err = s.db.QueryRow(checkAdminQuery, idUser).Scan(&existingAdminRestaurantId)
        
        if err == sql.ErrNoRows {
            // User doesn't have adminRestaurant record, create one
            insertQuery := `INSERT INTO adminRestaurant (idAdminRestaurant, idProfile) VALUES (?, ?)`
            _, err = s.db.Exec(insertQuery, newRoleId, idUser)
            if err != nil {
                return fmt.Errorf("error creating adminRestaurant role: %v", err)
            }
        } else if err != nil {
            return fmt.Errorf("error checking existing adminRestaurant: %v", err)
        }
        
        // Note: This method doesn't assign to specific restaurant, just creates the admin role
        // Use AssignUserToRoleWithEntity to assign to specific restaurant
        
    default:
        return fmt.Errorf("invalid role: %s", role)
    }
    
    return nil
}

// CreateNotification creates a new notification
func (s *Store) CreateNotification(notification types.NotificationCreation) (string, error) {
    // Generate a unique ID for the new notification
    newID, err := utils.CreateAnId()
    if err != nil {
        return "", fmt.Errorf("error generating ID for notification: %v", err)
    }

    query := `INSERT INTO notifications (idNotification, idAdmin, titre, type, description) VALUES (?, ?, ?, ?, ?)`
    _, err = s.db.Exec(query, newID, notification.IdAdmin, notification.Titre, notification.Type, notification.Description)
    if err != nil {
        return "", fmt.Errorf("error creating notification: %v", err)
    }

    return newID, nil
}

// GetNotificationsByAdmin retrieves all notifications for a specific admin
func (s *Store) GetNotificationsByAdmin(idAdmin string) ([]types.Notification, error) {
    query := `SELECT idNotification, idAdmin, titre, type, description FROM notifications WHERE idAdmin = ? ORDER BY idNotification DESC`
    rows, err := s.db.Query(query, idAdmin)
    if err != nil {
        return nil, fmt.Errorf("error getting notifications by admin: %v", err)
    }
    defer rows.Close()

    var notifications []types.Notification
    for rows.Next() {
        var notification types.Notification
        err := rows.Scan(
            &notification.IdNotification,
            &notification.IdAdmin,
            &notification.Titre,
            &notification.Type,
            &notification.Description,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning notification row: %v", err)
        }
        notifications = append(notifications, notification)
    }

    return notifications, nil
}

// GetAllNotifications retrieves all notifications
func (s *Store) GetAllNotifications() ([]types.Notification, error) {
    query := `SELECT idNotification, idAdmin, titre, type, description FROM notifications ORDER BY idNotification DESC`
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error getting all notifications: %v", err)
    }
    defer rows.Close()

    var notifications []types.Notification
    for rows.Next() {
        var notification types.Notification
        err := rows.Scan(
            &notification.IdNotification,
            &notification.IdAdmin,
            &notification.Titre,
            &notification.Type,
            &notification.Description,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning notification row: %v", err)
        }
        notifications = append(notifications, notification)
    }

    return notifications, nil
}

// CreateFeedback creates new feedback from a client
func (s *Store) CreateFeedback(feedback types.FeedbackCreation) error {
    query := `INSERT INTO feedback (idClient, comment) VALUES (?, ?)`
    _, err := s.db.Exec(query, feedback.IdClient, feedback.Comment)
    if err != nil {
        return fmt.Errorf("error creating feedback: %v", err)
    }

    return nil
}

// GetAllFeedbackWithClientInfo retrieves all feedback with client information
func (s *Store) GetAllFeedbackWithClientInfo() ([]types.Feedback, error) {
    query := `
        SELECT 
            f.idFeedback, f.idClient, f.comment, f.createdAt,
            p.firstName, p.lastName, p.email, c.username
        FROM feedback f
        JOIN client c ON f.idClient = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        ORDER BY f.createdAt DESC
    `
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error getting all feedback: %v", err)
    }
    defer rows.Close()

    var feedbacks []types.Feedback
    for rows.Next() {
        var feedback types.Feedback
        err := rows.Scan(
            &feedback.IdFeedback,
            &feedback.IdClient,
            &feedback.Comment,
            &feedback.CreatedAt,
            &feedback.ClientFirstName,
            &feedback.ClientLastName,
            &feedback.ClientEmail,
            &feedback.ClientUsername,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning feedback row: %v", err)
        }
        feedbacks = append(feedbacks, feedback)
    }

    return feedbacks, nil
}

func (s *Store) UpdateActivityAdmin(idActivity string, idAdminActivity string) error {
	// First check if activity exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM activity WHERE idActivity = ?)`
	err := s.db.QueryRow(checkQuery, idActivity).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking activity existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("activity not found")
	}

	// Update activity with admin (will replace existing admin if any)
	query := `UPDATE activity SET idAdminActivity = ? WHERE idActivity = ?`
	result, err := s.db.Exec(query, idAdminActivity, idActivity)
	if err != nil {
		return fmt.Errorf("error updating activity admin: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no activity found to update")
	}

	return nil
}

// AssignUserToRoleWithEntity assigns a user to a specific role with entity assignment
func (s *Store) AssignUserToRoleWithEntity(idUser string, role string, idActivity string, idRestaurant string) error {
    // First check if user exists
    var userExists bool
    checkQuery := `SELECT EXISTS(SELECT 1 FROM profile WHERE idProfile = ?)`
    err := s.db.QueryRow(checkQuery, idUser).Scan(&userExists)
    if err != nil {
        return fmt.Errorf("error checking user existence: %v", err)
    }
    if !userExists {
        return fmt.Errorf("user not found")
    }
    
    // Generate new ID for the role
    newRoleId, err := utils.CreateAnId()
    if err != nil {
        return fmt.Errorf("error generating role ID: %v", err)
    }
    
    // Assign the role based on type
    switch role {
    case "adminActivity":
        // Check if user is already ACTIVELY assigned to any activity as admin
        // This means they exist in adminActivity table AND are assigned to an activity
        var existingActiveActivityCount int
        checkActiveQuery := `
            SELECT COUNT(*) 
            FROM adminActivity aa
            JOIN activity a ON aa.idAdminActivity = a.idAdminActivity
            WHERE aa.idProfile = ?
        `
        err = s.db.QueryRow(checkActiveQuery, idUser).Scan(&existingActiveActivityCount)
        if err != nil {
            return fmt.Errorf("error checking existing active adminActivity assignment: %v", err)
        }
        if existingActiveActivityCount > 0 {
            return fmt.Errorf("user is already actively assigned as admin to an activity. An admin can only manage one activity")
        }
        
        // Check if the activity exists
        var activityExists bool
        checkActivityQuery := `SELECT EXISTS(SELECT 1 FROM activity WHERE idActivity = ?)`
        err = s.db.QueryRow(checkActivityQuery, idActivity).Scan(&activityExists)
        if err != nil {
            return fmt.Errorf("error checking activity existence: %v", err)
        }
        if !activityExists {
            return fmt.Errorf("activity not found")
        }
        
        // Check if user already has an adminActivity record, if not create one
        var existingAdminActivityId string
        checkAdminQuery := `SELECT idAdminActivity FROM adminActivity WHERE idProfile = ?`
        err = s.db.QueryRow(checkAdminQuery, idUser).Scan(&existingAdminActivityId)
        
        if err == sql.ErrNoRows {
            // User doesn't have adminActivity record, create one
            insertQuery := `INSERT INTO adminActivity (idAdminActivity, idProfile) VALUES (?, ?)`
            _, err = s.db.Exec(insertQuery, newRoleId, idUser)
            if err != nil {
                return fmt.Errorf("error creating adminActivity role: %v", err)
            }
            existingAdminActivityId = newRoleId
        } else if err != nil {
            return fmt.Errorf("error checking existing adminActivity: %v", err)
        }
        
        // Update activity with the admin (this will replace existing admin if any)
        err = s.UpdateActivityAdmin(idActivity, existingAdminActivityId)
        if err != nil {
            // If we created a new admin record, rollback
            if existingAdminActivityId == newRoleId {
                s.db.Exec(`DELETE FROM adminActivity WHERE idAdminActivity = ?`, newRoleId)
            }
            return fmt.Errorf("error assigning admin to activity: %v", err)
        }
        
    case "adminRestaurant":
        // Check if user is already ACTIVELY assigned to any restaurant as admin
        // This means they exist in adminRestaurant table AND are assigned to a restaurant
        var existingActiveRestaurantCount int
        checkActiveQuery := `
            SELECT COUNT(*) 
            FROM adminRestaurant ar
            JOIN restaurant r ON ar.idAdminRestaurant = r.idAdminRestaurant
            WHERE ar.idProfile = ?
        `
        err = s.db.QueryRow(checkActiveQuery, idUser).Scan(&existingActiveRestaurantCount)
        if err != nil {
            return fmt.Errorf("error checking existing active adminRestaurant assignment: %v", err)
        }
        if existingActiveRestaurantCount > 0 {
            return fmt.Errorf("user is already actively assigned as admin to a restaurant. An admin can only manage one restaurant")
        }
        
        // Check if the restaurant exists
        var restaurantExists bool
        checkRestaurantQuery := `SELECT EXISTS(SELECT 1 FROM restaurant WHERE idRestaurant = ?)`
        err = s.db.QueryRow(checkRestaurantQuery, idRestaurant).Scan(&restaurantExists)
        if err != nil {
            return fmt.Errorf("error checking restaurant existence: %v", err)
        }
        if !restaurantExists {
            return fmt.Errorf("restaurant not found")
        }
        
        // Check if user already has an adminRestaurant record, if not create one
        var existingAdminRestaurantId string
        checkAdminQuery := `SELECT idAdminRestaurant FROM adminRestaurant WHERE idProfile = ?`
        err = s.db.QueryRow(checkAdminQuery, idUser).Scan(&existingAdminRestaurantId)
        
        if err == sql.ErrNoRows {
            // User doesn't have adminRestaurant record, create one
            insertQuery := `INSERT INTO adminRestaurant (idAdminRestaurant, idProfile) VALUES (?, ?)`
            _, err = s.db.Exec(insertQuery, newRoleId, idUser)
            if err != nil {
                return fmt.Errorf("error creating adminRestaurant role: %v", err)
            }
            existingAdminRestaurantId = newRoleId
        } else if err != nil {
            return fmt.Errorf("error checking existing adminRestaurant: %v", err)
        }
        
        // Update restaurant with the admin (this will replace existing admin if any)
        err = s.UpdateRestaurantAdmin(idRestaurant, existingAdminRestaurantId)
        if err != nil {
            // If we created a new admin record, rollback
            if existingAdminRestaurantId == newRoleId {
                s.db.Exec(`DELETE FROM adminRestaurant WHERE idAdminRestaurant = ?`, newRoleId)
            }
            return fmt.Errorf("error assigning admin to restaurant: %v", err)
        }
        
    default:
        return fmt.Errorf("invalid role: %s", role)
    }
    
    return nil
}

func (s *Store) GetClientFollowersAndFollowing(idClient string) (*types.FollowListResponse, error) {
    response := &types.FollowListResponse{
        Following: []types.FollowingFollowerInfo{},
        Followers: []types.FollowingFollowerInfo{},
    }

    // Get following (users this client follows)
    followingQuery := `
        SELECT 
            c.idClient,
            c.username,
            p.firstName,
            p.lastName
        FROM friendship f
        JOIN client c ON f.idClient2 = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE f.idClient1 = ? AND f.status = 'accepted'
        ORDER BY p.firstName, p.lastName
    `
    
    rows, err := s.db.Query(followingQuery, idClient)
    if err != nil {
        return nil, fmt.Errorf("error retrieving following list: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var followingInfo types.FollowingFollowerInfo
        err := rows.Scan(
            &followingInfo.ClientId,
            &followingInfo.Username,
            &followingInfo.FirstName,
            &followingInfo.LastName,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning following row: %v", err)
        }
        response.Following = append(response.Following, followingInfo)
    }
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating following rows: %v", err)
    }

    // Get followers (users that follow this client)
    followersQuery := `
        SELECT 
            c.idClient,
            c.username,
            p.firstName,
            p.lastName
        FROM friendship f
        JOIN client c ON f.idClient1 = c.idClient
        JOIN profile p ON c.idProfile = p.idProfile
        WHERE f.idClient2 = ? AND f.status = 'accepted'
        ORDER BY p.firstName, p.lastName
    `
    
    rows, err = s.db.Query(followersQuery, idClient)
    if err != nil {
        return nil, fmt.Errorf("error retrieving followers list: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var followerInfo types.FollowingFollowerInfo
        err := rows.Scan(
            &followerInfo.ClientId,
            &followerInfo.Username,
            &followerInfo.FirstName,
            &followerInfo.LastName,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning follower row: %v", err)
        }
        response.Followers = append(response.Followers, followerInfo)
    }
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating follower rows: %v", err)
    }

    return response, nil
}

func (s *Store) CheckEmailExists(email string) (bool, error) {
    if email == "" {
        return false, nil
    }
    
    query := `SELECT EXISTS(SELECT 1 FROM profile WHERE email = ?)`
    var exists bool
    err := s.db.QueryRow(query, email).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("error checking email existence: %v", err)
    }
    return exists, nil
}

func (s *Store) CheckUsernameExists(username string) (bool, error) {
    if username == "" {
        return false, nil
    }
    
    query := `SELECT EXISTS(SELECT 1 FROM client WHERE username = ?)`
    var exists bool
    err := s.db.QueryRow(query, username).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("error checking username existence: %v", err)
    }
    return exists, nil
}

func (s *Store) CheckEmailAndUsernameAvailability(email, username string) (*types.AvailabilityCheckResponse, error) {
    response := &types.AvailabilityCheckResponse{
        EmailExists:    false,
        UsernameExists: false,
        Available:      true,
    }
    
    // Check email if provided
    if email != "" {
        emailExists, err := s.CheckEmailExists(email)
        if err != nil {
            return nil, err
        }
        response.EmailExists = emailExists
    }
    
    // Check username if provided
    if username != "" {
        usernameExists, err := s.CheckUsernameExists(username)
        if err != nil {
            return nil, err
        }
        response.UsernameExists = usernameExists
    }
    
    // Set availability - true only if neither exists (when provided)
    response.Available = !response.EmailExists && !response.UsernameExists
    
    return response, nil
}
