package middleware

import (
	"fmt"
	"net/http"
	"task-board/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AnonymousUserMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		anonymousUserID := c.GetHeader("X-Anonymous-User-Id")
		if anonymousUserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Anonymous-User-Id header required"})
			c.Abort()
			return
		}

		// Use UUID as unique identifier in email field
		userEmail := fmt.Sprintf("%s@anonymous.local", anonymousUserID)
		
		var user domain.User
		result := db.Where("email = ?", userEmail).First(&user)
		
		if result.Error == gorm.ErrRecordNotFound {
			// Create new anonymous user
			username := fmt.Sprintf("User_%s", anonymousUserID[:8])
			user = domain.User{
				Email:     userEmail,
				Username:  username,
				Password:  "anonymous", // Password not used for anonymous users
				FirstName: "Anonymous",
				LastName:  "User",
			}
			
			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create anonymous user"})
				c.Abort()
				return
			}
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Next()
	}
}

