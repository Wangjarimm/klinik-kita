package handlers

import (
	"net/http"
	"strings"

	"klinik/internal/database"
	"klinik/internal/middleware"
	"klinik/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerInput struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func Register(c *gin.Context) {
	var in registerInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	var count int64
	database.DB.Model(&models.User{}).Where("username = ? OR email = ?", in.Username, in.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username atau email sudah terdaftar"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}

	user := models.User{Username: in.Username, Email: in.Email, Password: string(hash), Role: models.RolePasien}
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return tx.Create(&models.Pasien{UserID: &user.ID, NamaLengkap: in.Username}).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registrasi berhasil", "role": user.Role})
}

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var in loginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", strings.TrimSpace(in.Username)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	middleware.SetSessionCookie(c, jwtSecret, middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "username": user.Username, "role": user.Role})
}

func Logout(c *gin.Context) {
	middleware.ClearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Keluar berhasil"})
}

func Me(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	c.JSON(http.StatusOK, gin.H{"user_id": claims.UserID, "username": claims.Username, "role": claims.Role})
}
