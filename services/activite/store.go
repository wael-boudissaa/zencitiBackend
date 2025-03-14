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

func (s *Store) GetActivite() (*[]types.Activite, error) {
    query := `SELECT * FROM activite`
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close() // Ensure rows are closed to avoid memory leaks
    var activite []types.Activite

    for rows.Next() {
        var act types.Activite
        err = rows.Scan(
            &act.IdActivite,
            &act.NameActivite,
            &act.Description,
        )
        if err != nil {
            return nil, err
        }
        activite = append(activite, act)
    }
    if err:=rows.Err(); err!=nil {
        return nil, err
    }
    return &activite, nil
}


func (s *Store) GetActiviteById(id int) (*types.Activite, error) {
    query := `SELECT * FROM activite WHERE idActivite = ?`
    row := s.db.QueryRow(query, id)
    var act types.Activite
    err := row.Scan(
        &act.IdActivite,
        &act.NameActivite,
        &act.Description,
    )
    if err != nil {
        return nil, err
    }
    return &act, nil
}

func (s *Store) GetActiviteTypes() (*[]types.ActivitetType, error) {
    query := `SELECT * FROM activiteType`
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
        )
        if err != nil {
            return nil, err
        }
        activite = append(activite, act)
    }
    if err:=rows.Err(); err!=nil {
        return nil, err
    }
    return &activite, nil
}


