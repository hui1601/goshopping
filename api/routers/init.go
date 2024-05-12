package routers

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"math/rand"
	"net/http"
	"os"
)

var db *sql.DB

const serverErrorMessage = "The Quick Brown Fox Failed To Jump Over The Lazy Dog"

func init() {
	gob.Register(gin.H{})
	gob.Register(purchaseRequest{})
	connectDB()
	createTable()
}

func createTable() {
	createImageTable()
	createUserTable()
	createProductsTable()
	createPurchaseTables()
}

func RegisterRoutes(router *gin.Engine) {
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	router.Use(sessions.Sessions("session", store))
	RegisterAdminRoutes(router)
	RegisterImageRoutes(router)
	RegisterAuthRoutes(router)
	RegisterUserRoutes(router)
	RegisterProductsRoutes(router)
	RegisterPurchaseRoutes(router)
}
func connectDB() {
	// Connect to database
	println(os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_DATABASE"))
	dbConn, err := sql.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_DATABASE"))
	if err != nil {
		panic(err.Error())
	}
	db = dbConn
}

func createRandomID() string {
	stringBytes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
	randomID := make([]byte, 12)
	for i := range randomID {
		randomID[i] = stringBytes[rand.Intn(len(stringBytes))]
	}
	return string(randomID)
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !checkUserID(user.(string)) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func hashPassword(password string) string {
	hash := sha512.New()
	_, _ = hash.Write([]byte(password))
	return encodeBase64(hash.Sum(nil))
}

func AdminRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user").(string)
	if getUserType(user) != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getUserType(id string) string {
	stmt, err := db.Prepare("SELECT user_type FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var userType string
	err = stmt.QueryRow(id).Scan(&userType)
	if err != nil {
		panic(err.Error())
	}
	return userType
}
