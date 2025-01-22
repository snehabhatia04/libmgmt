package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snehabhatia04/libmgmt/model" // Import the model package
)

// Handler for getting all users
func GetUsers(c *gin.Context) {
    db, err := sql.Open("postgres", "user=postgres dbname=demoDB password=Sneha host=localhost port=5432 sslmode=disable")
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
        return
    }
    defer db.Close()

    var users []model.User
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var user model.User
        if err := rows.Scan(&user.ID, &user.Name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        users = append(users, user)
    }

    c.JSON(http.StatusOK, users)
}

// Handler for adding a new user
func AddUser(c *gin.Context) {
    var user model.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db, err := sql.Open("postgres", "user=postgres dbname=demoDB password=Sneha host=localhost port=5432 sslmode=disable")
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
        return
    }
    defer db.Close()

    // Insert new user into the database
    _, err = db.Exec("INSERT INTO users (name) VALUES ($1)", user.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User added"})
}

// Handler for adding a book
func AddBook(c *gin.Context) {
    var book model.Book
    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db, err := sql.Open("postgres", "user=postgres dbname=demoDB password=Sneha host=localhost port=5432 sslmode=disable")
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
        return
    }
    defer db.Close()

    // Insert the book into the database
    _, err = db.Exec("INSERT INTO books (title, author_id, location) VALUES ($1, $2, $3)", book.Title, book.AuthorID, book.Location)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Book added"})
}

// Handler for issuing a book to a user
func IssueBook(c *gin.Context) {
    var bookIssue model.BookIssue
    if err := c.ShouldBindJSON(&bookIssue); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Calculate the due date (15 days from issue date)
    issueDate, err := time.Parse("2006-01-02", bookIssue.IssueDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue date format"})
        return
    }
    dueDate := issueDate.Add(15 * 24 * time.Hour).Format("2006-01-02")

    db, err := sql.Open("postgres", "user=postgres dbname=demoDB password=Sneha host=localhost port=5432 sslmode=disable")
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
        return
    }
    defer db.Close()

    // Insert the issue record into the database
    _, err = db.Exec("INSERT INTO book_issues (user_id, book_id, issue_date, due_date) VALUES ($1, $2, $3, $4)", bookIssue.UserID, bookIssue.BookID, bookIssue.IssueDate, dueDate)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Book issued successfully"})
}

// Handler for returning a book and calculating the fine
func ReturnBook(c *gin.Context) {
    var bookIssue model.BookIssue
    if err := c.ShouldBindJSON(&bookIssue); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db, err := sql.Open("postgres", "user=postgres dbname=demoDB password=Sneha host=localhost port=5432 sslmode=disable")
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
        return
    }
    defer db.Close()

    // Check if the book is issued to the user
    var existingBookIssue model.BookIssue
    err = db.QueryRow("SELECT * FROM book_issues WHERE user_id = $1 AND book_id = $2", bookIssue.UserID, bookIssue.BookID).Scan(&existingBookIssue.ID, &existingBookIssue.UserID, &existingBookIssue.BookID, &existingBookIssue.IssueDate, &existingBookIssue.ReturnDate, &existingBookIssue.DueDate)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not issued to this user"})
        return
    }

    // Parse the return date
    returnDate, err := time.Parse("2006-01-02", bookIssue.ReturnDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return date format"})
        return
    }

    // Calculate fine if the return date is after the due date
    dueDate, err := time.Parse("2006-01-02", existingBookIssue.DueDate)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing due date"})
        return
    }

    fineAmount := 0.0
    if returnDate.After(dueDate) {
        daysLate := int(returnDate.Sub(dueDate).Hours() / 24)
        fineAmount = float64(daysLate) * 10.0 // 10 units per day
    }

    // Insert fine record
    _, err = db.Exec("INSERT INTO fines (user_id, fine_amount, date) VALUES ($1, $2, NOW())", bookIssue.UserID, fineAmount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting fine"})
        return
    }

    // Update the book return date
    _, err = db.Exec("UPDATE book_issues SET return_date = $1 WHERE user_id = $2 AND book_id = $3", bookIssue.ReturnDate, bookIssue.UserID, bookIssue.BookID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully", "fine": fineAmount})
}
