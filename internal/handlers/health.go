package handlers

import (
	"net/http"

	"e-commerce/internal/db"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
