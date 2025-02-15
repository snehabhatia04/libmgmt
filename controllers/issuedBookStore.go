package controllers

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/snehabhatia04/libmgmt/model"
)

type IssuedBookStore struct {
	DB *sqlx.DB
}

// NewDBIssuedBookStore initializes a new IssuedBookStore
func NewDBIssuedBookStore(db *sqlx.DB) *IssuedBookStore {
	return &IssuedBookStore{DB: db}
}

// IssueBook issues a book to a user
func (store *IssuedBookStore) IssueBook(issue *model.BookIssue) error {
	query := `INSERT INTO book_issues (user_id, book_id, issue_date, due_date) VALUES ($1, $2, $3, $4) RETURNING id`
	err := store.DB.QueryRow(query, issue.UserID, issue.BookID, issue.IssueDate, issue.DueDate).Scan(&issue.ID)
	if err != nil {
		log.Println("Error issuing book:", err)
		return err
	}
	return nil
}

// ReturnBook processes the return of a book and calculates fine if applicable
func (store *IssuedBookStore) ReturnBook(bookID int) (float64, error) {
	var issue model.BookIssue
	query := `SELECT * FROM book_issues WHERE book_id = $1`
	err := store.DB.Get(&issue, query, bookID)
	if err != nil {
		log.Println("Book not issued or not found:", err)
		return 0, errors.New("book not issued or not found")
	}

	// Calculate fine if returned late
	dueDate, _ := time.Parse("2006-01-02", issue.DueDate)
	returnDate := time.Now()
	var fineAmount float64
	if returnDate.After(dueDate) {
		daysLate := returnDate.Sub(dueDate).Hours() / 24
		fineAmount = daysLate * 5 // Example: ₹5 fine per late day
	}

	// Delete issued record
	delQuery := `DELETE FROM book_issues WHERE book_id = $1`
	_, err = store.DB.Exec(delQuery, bookID)
	if err != nil {
		log.Println("Error deleting issued record:", err)
		return fineAmount, err
	}

	return fineAmount, nil
}

// GetIssuedBookByBookID retrieves issued book details by Book ID
func (store *IssuedBookStore) GetIssuedBookByBookID(bookID int) (model.BookIssue, error) {
	var issue model.BookIssue
	query := `SELECT * FROM book_issues WHERE book_id = $1`
	err := store.DB.Get(&issue, query, bookID)
	if err != nil {
		log.Println("Issued book not found:", err)
		return model.BookIssue{}, errors.New("issued book not found")
	}
	return issue, nil
}

// GetIssuedBooks retrieves all issued books
func (store *IssuedBookStore) GetIssuedBooks() ([]model.BookIssue, error) {
	var issues []model.BookIssue
	query := `SELECT * FROM book_issues`
	err := store.DB.Select(&issues, query)
	if err != nil {
		log.Println("Error fetching issued books:", err)
		return nil, err
	}
	return issues, nil
}
