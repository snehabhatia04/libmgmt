//

package controllers

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/snehabhatia04/libmgmt/model"
)

type BookStore struct {
	DB *sqlx.DB
}

// NewDBBookStore initializes a new BookStore
func NewDBBookStore(db *sqlx.DB) *BookStore {
	return &BookStore{DB: db}
}

// // GetBooks retrieves all books
// func (store *BookStore) GetBooks() ([]model.Book, error) {
// 	var books []model.Book
// 	query := `SELECT * FROM books`
// 	err := store.DB.Select(&books, query)
// 	if err != nil {
// 		log.Println("Error fetching books:", err)
// 		return nil, err
// 	}
// 	return books, nil
// }

// // GetBook retrieves a specific book by ID
// func (store *BookStore) GetBook(id int) (*model.Book, error) {
// 	var book model.Book
// 	query := `SELECT * FROM books WHERE id = $1`
// 	err := store.DB.Get(&book, query, id)
// 	if err != nil {
// 		return nil, errors.New("book not found")
// 	}
// 	return &book, nil
// }


// GetBooks retrieves all books
func (store *BookStore) GetBooks() ([]model.Book, error) {
	var books []model.Book
	query := `SELECT * FROM books`
	err := store.DB.Select(&books, query)
	if err != nil {
		log.Println("Error fetching books:", err)
		return nil, err
	}
	return books, nil
}

// FIX: GetBook should return (model.Book, error) instead of (*model.Book, error)
func (store *BookStore) GetBook(id int) (model.Book, error) {
	var book model.Book
	query := `SELECT * FROM books WHERE id = $1`
	err := store.DB.Get(&book, query, id)
	if err != nil {
		log.Println("Book not found:", err)
		return model.Book{}, errors.New("book not found") // Return empty struct instead of nil
	}
	return book, nil
}


// CreateBook adds a new book
func (store *BookStore) CreateBook(book *model.Book) error {
	query := `INSERT INTO books (title, author_id, location) VALUES ($1, $2, $3) RETURNING id`
	err := store.DB.QueryRow(query, book.Title, book.AuthorID, book.Location).Scan(&book.ID)
	if err != nil {
		log.Println("Error adding book:", err)
		return err
	}
	return nil
}

// UpdateBook updates an existing book
func (store *BookStore) UpdateBook(book *model.Book) error {
	query := `UPDATE books SET title = $1, author_id = $2, location = $3 WHERE id = $4`
	_, err := store.DB.Exec(query, book.Title, book.AuthorID, book.Location, book.ID)
	if err != nil {
		log.Println("Error updating book:", err)
		return err
	}
	return nil
}

// DeleteBook deletes a book by ID
func (store *BookStore) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := store.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting book:", err)
		return err
	}
	return nil
}
