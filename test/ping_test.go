package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Aircraft API is running...",
			"version": "1.0.0",
		})
	})

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)                            
	assert.Contains(t, w.Body.String(), "Aircraft API is running...") 
	assert.Contains(t, w.Body.String(), "1.0.0")                      
}
