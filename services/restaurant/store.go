package restaurant

import (
	"database/sql"

	"github.com/wael-boudissaa/zencitiBackend/types"
)

type store struct {
    db *sql.DB
}

func NewStore(db *sql.DB) *store {
    return &store{db: db}
}

func (s *store) GetRestaurant() (*[]types.Restaurant, error) {
    query := `SELECT * FROM restaurant`
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close() // Ensure rows are closed to avoid memory leaks
    var restaurant []types.Restaurant

    for rows.Next() {
        var rest types.Restaurant
        err = rows.Scan(
            &rest.IdRestaurant,
            &rest.NameRestaurant,
            &rest.Description,
            &rest.IdActivite,
            &rest.Image,
            &rest.Location,
            &rest.Capacity,
        )
        if err != nil {
            return nil, err
        }
        restaurant = append(restaurant, rest)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return &restaurant, nil
}

func (s *store) GetRestaurantById(id string) (*types.Restaurant, error) {
    query := `SELECT * FROM restaurant WHERE idRestaurant = ?`
    row := s.db.QueryRow(query, id)
    var rest types.Restaurant
    err := row.Scan(
        &rest.IdRestaurant,
        &rest.NameRestaurant,
        &rest.Description,
        &rest.IdActivite,
        &rest.Image,
        &rest.Location,
        &rest.Capacity,
    )
    if err != nil {
        return nil, err
    }
    return &rest, nil
}


func (s *store) GetRestaurantWorkers ()
