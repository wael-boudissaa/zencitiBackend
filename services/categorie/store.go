package categorie

import (
	"database/sql"

	"github.com/wael-boudissaa/marquinoBackend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db}
}

func (s *Store) GetCategories() (*[]types.Categorie, error) {
	query := `Select * from categorie`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	cateogireResult, err := scanRowsIntoCategorie(rows)
	if err != nil {
		return nil, err
	}
	return cateogireResult, nil
}

func scanRowsIntoCategorie(rows *sql.Rows) (*[]types.Categorie, error) {
	categories := []types.Categorie{} // Initialize an empty slice (no pointer needed)
	for rows.Next() {
		categorie := types.Categorie{} // Create a new instance inside the loop
		err := rows.Scan(&categorie.IdCategorie, &categorie.NameCategorie)
		if err != nil {
			return nil, err
		}
		categories = append(categories, categorie)
	}
	return &categories, nil
}
