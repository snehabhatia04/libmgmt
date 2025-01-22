package controllers

import (
	"net/http"
	"your_project/database"

	"github.com/gin-gonic/gin"
	"main.go/model"
)

// Create Book Issue
func CreateBookIssue(c *gin.Context) {
	var bookIssue model.BookIssue
	if err := c.ShouldBindJSON(&bookIssue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := database.GetDB()
	if err := db.Create(&bookIssue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookIssue)
}

// Get all book issues
func GetBookIssues(c *gin.Context) {
	var bookIssues []model.BookIssue
	db := database.GetDB()
	if err := db.Find(&bookIssues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookIssues)
}

// Get book issue by ID
func GetBookIssue(c *gin.Context) {
	id := c.Param("id")
	var bookIssue model.BookIssue
	db := database.GetDB()
	if err := db.First(&bookIssue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book Issue not found"})
		return
	}
	c.JSON(http.StatusOK, bookIssue)
}

// Update book issue
func UpdateBookIssue(c *gin.Context) {
	id := c.Param("id")
	var bookIssue model.BookIssue
	if err := c.ShouldBindJSON(&bookIssue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := database.GetDB()
	if err := db.Model(&bookIssue).Where("id = ?", id).Updates(bookIssue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookIssue)
}

// Delete book issue
func DeleteBookIssue(c *gin.Context) {
	id := c.Param("id")
	var bookIssue model.BookIssue
	db := database.GetDB()
	if err := db.Delete(&bookIssue, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book Issue deleted successfully"})
}
