package controllers

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/snehabhatia04/libmgmt/model"
)

type FineStore struct {
	DB *sqlx.DB
}

// NewDBFineStore initializes a new FineStore
func NewDBFineStore(db *sqlx.DB) *FineStore {
	return &FineStore{DB: db}
}

// AddFine adds a new fine record
func (store *FineStore) AddFine(fine *model.Fine) error {
	query := `INSERT INTO fines (user_id, fine_amount, date) VALUES ($1, $2, $3) RETURNING id`
	err := store.DB.QueryRow(query, fine.UserID, fine.FineAmount, fine.Date).Scan(&fine.ID)
	if err != nil {
		log.Println("Error adding fine:", err)
		return err
	}
	return nil
}

// GetFine retrieves a single fine by ID
func (store *FineStore) GetFine(id int) (model.Fine, error) {
	var fine model.Fine
	query := `SELECT * FROM fines WHERE id = $1`
	err := store.DB.Get(&fine, query, id)
	if err != nil {
		log.Println("Fine not found:", err)
		return model.Fine{}, errors.New("fine not found")
	}
	return fine, nil
}

// GetFines retrieves all fines
func (store *FineStore) GetFines() ([]model.Fine, error) {
	var fines []model.Fine
	query := `SELECT * FROM fines`
	err := store.DB.Select(&fines, query)
	if err != nil {
		log.Println("Error fetching fines:", err)
		return nil, err
	}
	return fines, nil
}

// DeleteFine removes a fine record by ID
func (store *FineStore) DeleteFine(id int) error {
	query := `DELETE FROM fines WHERE id = $1`
	result, err := store.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting fine:", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("fine not found")
	}
	return nil
}


// package controllers

// import (
// 	"log"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/gin-gonic/gin"
// 	"github.com/snehabhatia04/libmgmt/model"
// )

// type FineStore struct {
// 	DB *sqlx.DB
// }

// // NewDBFineStore initializes a new FineStore
// func NewDBFineStore(db *sqlx.DB) *FineStore {
// 	return &FineStore{DB: db}
// }

// // GetFines retrieves all fines
// func (store *FineStore) GetFines(c *gin.Context) {
// 	var fines []model.Fine
// 	query := `SELECT * FROM fines`
// 	err := store.DB.Select(&fines, query)
// 	if err != nil {
// 		log.Println("Error fetching fines:", err)
// 		c.JSON(500, gin.H{"error": "Failed to retrieve fines"})
// 		return
// 	}
// 	c.JSON(200, fines)
// }

// // GetFine retrieves a specific fine by ID
// func (store *FineStore) GetFine(c *gin.Context) {
// 	fineID := c.Param("id")
// 	var fine model.Fine

// 	query := `SELECT * FROM fines WHERE id = $1`
// 	err := store.DB.Get(&fine, query, fineID)
// 	if err != nil {
// 		c.JSON(404, gin.H{"error": "Fine not found"})
// 		return
// 	}
// 	c.JSON(200, fine)
// }

// // DeleteFine deletes a fine by ID
// func (store *FineStore) DeleteFine(c *gin.Context) {
// 	fineID := c.Param("id")
// 	query := `DELETE FROM fines WHERE id = $1`
// 	_, err := store.DB.Exec(query, fineID)
// 	if err != nil {
// 		log.Println("Error deleting fine:", err)
// 		c.JSON(500, gin.H{"error": "Failed to delete fine"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "Fine deleted successfully"})
// }
