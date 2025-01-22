package controllers

import (
	"net/http"
	"your_project/database"
	"your_project/models"

	"github.com/gin-gonic/gin"
)

// Create Fine
func CreateFine(c *gin.Context) {
	var fine models.Fine
	if err := c.ShouldBindJSON(&fine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := database.GetDB()
	if err := db.Create(&fine).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fine)
}

// Get all fines
func GetFines(c *gin.Context) {
	var fines []models.Fine
	db := database.GetDB()
	if err := db.Find(&fines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fines)
}

// Get fine by ID
func GetFine(c *gin.Context) {
	id := c.Param("id")
	var fine models.Fine
	db := database.GetDB()
	if err := db.First(&fine, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fine not found"})
		return
	}
	c.JSON(http.StatusOK, fine)
}

// Update fine
func UpdateFine(c *gin.Context) {
	id := c.Param("id")
	var fine models.Fine
	if err := c.ShouldBindJSON(&fine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := database.GetDB()
	if err := db.Model(&fine).Where("id = ?", id).Updates(fine).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fine)
}

// Delete fine
func DeleteFine(c *gin.Context) {
	id := c.Param("id")
	var fine models.Fine
	db := database.GetDB()
	if err := db.Delete(&fine, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fine deleted successfully"})
}
