package commande

import (
	"database/sql"
	"fmt"

	"github.com/wael-boudissaa/marquinoBackend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

type Commande = types.Commande
type CommandeCreate = types.CommandeCreate

func (s *Store) GetAllCommandes() (*[]Commande, error) {
	query := `select * from commande`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	var commandes = new([]Commande)
	commandes, err = scanRowsIntoCommande(rows)
	if err != nil {
		return nil, err
	}
	return commandes, nil
}

func (s *Store) GetCommandeById(id string) (*Commande, error) {
	query := `select * from commande where id = ? `
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	var commande = new(Commande)
	err = rows.Scan(&commande.IdCommande, &commande.IdCustomer)
	if err != nil {
		return nil, err
	}
	return commande, nil
}

func (s *Store) CreateCommande(idCommande, idCustomer string, price int) error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `INSERT INTO commande (idCommande, idCustomer, price,createdAt ) VALUES (?, ?, ?,now()) )`
	result, err := s.db.Exec(query, idCommande, idCustomer, price)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting rows affected: %v\n", err)
		return err
	}
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	return nil
}

func (s *Store) InsertProductINCommande(product types.ProductBought, idCommande string) (*types.CommandeProduct, error) {
	query := `INSERT INTO commande_products (idCommande, idProduct, quantity) VALUES (?, ?, ?)`

	// Execute the insert query
	result, err := s.db.Exec(query, idCommande, product.IdProduct, product.Quantity)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	// Optionally, check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no rows were affected by the insert")
	}

	// Return the inserted values if needed
	commandeProduct := &types.CommandeProduct{
		IdCommande: idCommande,
		IdProduct:  product.IdProduct,
	}

	return commandeProduct, nil
}

// UpdateCommande(commande Commande) error
// DeleteCommande(commande Commande) error

func GetCommandeByUser(idUser string) (*[]Commande, error) {
	return nil, nil
}

func scanRowsIntoCommande(rows *sql.Rows) (*[]Commande, error) {
	commandes := new([]types.Commande)

	commande := new(types.Commande)
	for rows.Next() {
		err := rows.Scan(&commande.IdCommande, &commande.IdCustomer)
		if err != nil {
			return nil, err
		}
		*commandes = append(*commandes, *commande)

	}
	return nil, nil
}
