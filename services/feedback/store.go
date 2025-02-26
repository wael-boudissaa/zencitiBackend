package feedback

import (
	"database/sql"

	"github.com/wael-boudissaa/marquinoBackend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateFeedBack(idFeedBack, idCustomer, comment string) error {
	_, err := s.db.Exec("INSERT INTO feedback(idCustomer,idFeedback,comment,createdAt) VALUES(?,?,?,?)", idCustomer, idFeedBack, comment, `now()`)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAllFeedBack() (*[]types.FeedBack, error) {
	rows, err := s.db.Query("SELECT * FROM feedback")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedBacks []types.FeedBack
	for rows.Next() {
		var feedBack types.FeedBack
		err := rows.Scan(&feedBack.IDCustomer, &feedBack.IDFeedBack, &feedBack.Comment, &feedBack.CreatedAt)
		if err != nil {
			return nil, err
		}
		feedBacks = append(feedBacks, feedBack)
	}
	return &feedBacks, nil
}
