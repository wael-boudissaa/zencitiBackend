package activite

import (
	"database/sql"

	"github.com/wael-boudissaa/zencitiBackend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// func (s *Store) GetActivite() (*[]types.Activite, error) {
// 	query := `SELECT * FROM activite`
// 	rows, err := s.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close() // Ensure rows are closed to avoid memory leaks
// 	var activite []types.Activite
//
// 	for rows.Next() {
// 		var act types.Activite
// 		err = rows.Scan(
// 			&act.IdActivite,
// 			&act.NameActivite,
// 			&act.Description,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		activite = append(activite, act)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return &activite, nil
// }
//

func (s *Store) GetRecentActivities(idClient string) (*[]types.ActivityProfile, error) {
	query := `SELECT
    activity.idActivity,activity.nameActivity,activity.descriptionActivity,activity.imageActivity,activity.popularity,clientActivity.timeActivity
    FROM clientActivity join activity on clientActivity.idActivity=
    activity.idActivity where clientActivity.idClient=? ORDER BY
    clientActivity.timeActivity DESC LIMIT 5 `

	rows, err := s.db.Query(query, idClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed to avoid memory leaks
	var activite []types.ActivityProfile
	for rows.Next() {
		var act types.ActivityProfile
		err = rows.Scan(
			&act.IdActivity,
			&act.NameActivity,
			&act.Description,
			&act.ImageActivite,
			&act.Popularity,
			&act.TimeActivity,
		)
		if err != nil {
			return nil, err
		}
		activite = append(activite, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activite, nil
}

func (s *Store) GetActiviteById(id string) (*types.Activity, error) {
	query := `SELECT * FROM activity WHERE idActivity = ?`
	row := s.db.QueryRow(query, id)
	var act types.Activity
	err := row.Scan(
		&act.IdActivity,
		&act.NameActivity,
		&act.Description,
		&act.ImageActivite,
		&act.IdTypeActivity,
		&act.Popularity,
	)
	if err != nil {
		return nil, err
	}
	return &act, nil
}

func (s *Store) GetActiviteTypes() (*[]types.ActivitetType, error) {
	query := `SELECT * FROM typeActivity`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activite []types.ActivitetType

	for rows.Next() {
		var act types.ActivitetType
		err = rows.Scan(
			&act.IdActiviteType,
			&act.NameActiviteType,
			&act.ImageActivity,
		)
		if err != nil {
			return nil, err
		}
		activite = append(activite, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activite, nil
}

func (s *Store) GetActivityByTypes(id string) (*[]types.Activity, error) {
	query := `SELECT * FROM activity WHERE idTypeActivity = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activity []types.Activity

	for rows.Next() {
		var act types.Activity
		err = rows.Scan(
			&act.IdActivity,
			&act.NameActivity,
			&act.Description,
			&act.ImageActivite,
			&act.IdTypeActivity,
			&act.Popularity,
		)
		if err != nil {
			return nil, err
		}
		activity = append(activity, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *Store) GetPopularActivities() (*[]types.Activity, error) {
	query := `SELECT * FROM activity ORDER BY popularity DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activite []types.Activity

	for rows.Next() {
		var act types.Activity
		err = rows.Scan(
			&act.IdActivity,
			&act.NameActivity,
			&act.Description,
			&act.ImageActivite,
			&act.IdTypeActivity,
			&act.Popularity,
		)
		if err != nil {
			return nil, err
		}
		activite = append(activite, act)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &activite, nil
}
