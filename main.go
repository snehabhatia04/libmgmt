package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// User struct
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Book struct
type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
	Location string `json:"location"`
}

// BookIssue struct
type BookIssue struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	BookID     int    `json:"book_id"`
	IssueDate  string `json:"issue_date"`
	ReturnDate string `json:"return_date,omitempty"`
	DueDate    string `json:"due_date"`
}

// Fine struct
type Fine struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	FineAmount float64 `json:"fine_amount"`
	Date       string  `json:"date"`
}

// Connect to the PostgreSQL database
func initDB() {
    var err error
	db, err = sqlx.Connect("postgres", "user=postgres dbname=demoDB password=Sneha host=localhost port=5432 sslmode=disable")
    if err != nil {
        log.Fatalln("Failed to connect to database:", err)
    }
    fmt.Println("Database connected successfully.")
}


// Get all users
func getUsers(c *gin.Context) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

// Add a user
func addUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := db.NamedExec("INSERT INTO users (name) VALUES (:name)", user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "User added"})
}

// Get books by author
func getBooksByAuthor(c *gin.Context) {
	authorID := c.Param("author_id")
	var books []Book
	err := db.Select(&books, "SELECT * FROM books WHERE author_id = $1", authorID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, books)
}

// Add a book
func addBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := db.NamedExec("INSERT INTO books (title, author_id, location) VALUES (:title, :author_id, :location)", book)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Book added"})
}

// Issue a book to a user
func issueBook(c *gin.Context) {
	var issue BookIssue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Calculate the due date (15 days from issue date)
	issueDate, err := time.Parse("2006-01-02", issue.IssueDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid issue date format"})
		return
	}
	dueDate := issueDate.Add(15 * 24 * time.Hour).Format("2006-01-02")

	// Insert into book_issues table
	_, err = db.NamedExec("INSERT INTO book_issues (user_id, book_id, issue_date, due_date) VALUES (:user_id, :book_id, :issue_date, :due_date)", 
		map[string]interface{}{
			"user_id":  issue.UserID,
			"book_id":  issue.BookID,
			"issue_date": issue.IssueDate,
			"due_date": dueDate,
		})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Book issued successfully"})
}

// Return a book and calculate fine if applicable
func returnBook(c *gin.Context) {
	var issue BookIssue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check the issue details from the database
	var bookIssue BookIssue
	err := db.Get(&bookIssue, "SELECT * FROM book_issues WHERE user_id = $1 AND book_id = $2", issue.UserID, issue.BookID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Book not issued to this user"})
		return
	}

	// Parse return date and due date
	returnDate, err := time.Parse("2006-01-02", issue.ReturnDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid return date format"})
		return
	}

	dueDate, err := time.Parse("2006-01-02", bookIssue.DueDate)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error parsing due date"})
		return
	}

	// Calculate fine if return date is later than due date
	var fineAmount float64
	if returnDate.After(dueDate) {
		daysLate := int(returnDate.Sub(dueDate).Hours() / 24)
		fineAmount = float64(daysLate) * 10.0 // 10 units per day
		// Insert fine record
		_, err = db.NamedExec("INSERT INTO fines (user_id, fine_amount, date) VALUES (:user_id, :fine_amount, NOW())",
			map[string]interface{}{
				"user_id":    issue.UserID,
				"fine_amount": fineAmount,
			})
		if err != nil {
			c.JSON(500, gin.H{"error": "Error inserting fine"})
			return
		}
	}

	// Update the return date in book_issues table
	_, err = db.NamedExec("UPDATE book_issues SET return_date = :return_date WHERE user_id = :user_id AND book_id = :book_id",
		map[string]interface{}{
			"return_date": issue.ReturnDate,
			"user_id":     issue.UserID,
			"book_id":     issue.BookID,
		})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Book returned successfully", "fine": fineAmount})
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// Routes for users
	r.GET("/users", getUsers)
	r.POST("/users", addUser)

	// Routes for books
	r.GET("/authors/:author_id/books", getBooksByAuthor)
	r.POST("/books", addBook)

	// Routes for book issues and returns
	r.POST("/issue_book", issueBook)
	r.POST("/return_book", returnBook)

	r.Run(":8080")
}
