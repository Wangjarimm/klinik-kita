package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const CookieName = "session"

type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"usr"`
	Role     string `json:"rol"`
	jwt.RegisteredClaims
}

// shortcut: Secure=false karena seluruh stack demo berjalan di localhost http; aktifkan saat deploy https
func SetSessionCookie(c *gin.Context, secret []byte, cl Claims) {
	cl.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	signed, err := token.SignedString(secret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat sesi"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CookieName,
		Value:    signed,
		Path:     "/",
		MaxAge:   12 * 3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ParseToken(secret []byte, value string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(value, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	return claims, err
}

func AuthRequired(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(CookieName)
		if err != nil || cookie == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "belum login"})
			return
		}
		claims, err := ParseToken(secret, cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sesi tidak valid"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func RoleRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(*Claims)
		if claims.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "akses ditolak"})
			return
		}
		c.Next()
	}
}

func CORS(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Data kuota/antrean tidak boleh disimpan cache browser.
		c.Header("Cache-Control", "no-store")
		if c.GetHeader("Origin") == origin {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Logout menghapus cookie sesi.
func ClearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
