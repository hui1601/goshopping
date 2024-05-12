package routers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type Image struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

func RegisterImageRoutes(router *gin.Engine) {
	authorized := router.Group("/image", AuthRequired)
	admin := authorized.Group("", AdminRequired)
	admin.POST("", postImage)
	authorized.GET("/:id", getImage)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func checkImageTable() bool {
	// Check if image table exists
	stmt, err := db.Prepare("SHOW TABLES LIKE 'images'")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var table string
	err = stmt.QueryRow().Scan(&table)
	return err == nil
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func createImageTable() {
	if checkImageTable() {
		return
	}
	// create table
	stmt, err := db.Prepare("CREATE TABLE images (id VARCHAR(12) PRIMARY KEY, image LONGBLOB)")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec()
	if err != nil {
		panic(err.Error())
	}
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getImage(c *gin.Context) {
	// Get image by id
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"message": "ID is required",
		})
		return
	}
	stmt, err := db.Prepare("SELECT * FROM images WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var imageSt Image
	var imageBytes []byte
	err = stmt.QueryRow(id).Scan(&imageSt.ID, &imageBytes)
	if err != nil {
		log.Println(err.Error())
		c.JSON(404, gin.H{
			"message": "Image not found",
		})
		return
	}
	// encode imageSt to base64
	imageSt.Image = base64.StdEncoding.EncodeToString(imageBytes)
	imageSt.Image = "data:image/jpeg;base64," + imageSt.Image
	c.JSON(200, imageSt)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func postImage(c *gin.Context) {
	// Post imageSt
	var imageSt Image
	err := c.BindJSON(&imageSt)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid JSON",
		})
		return
	}
	if imageSt.Image == "" {
		c.JSON(400, gin.H{
			"message": "Image is required",
		})
		return
	}
	// decode imageSt from base64
	imageSt.Image = imageSt.Image[strings.IndexByte(imageSt.Image, ',')+1:]
	imageBytes, err := base64.StdEncoding.DecodeString(imageSt.Image)
	if err != nil || len(imageBytes) == 0 || http.DetectContentType(imageBytes) != "image/jpeg" {
		c.JSON(400, gin.H{
			"message":    "Invalid image format. Only jpeg is allowed",
			"postedType": http.DetectContentType(imageBytes),
		})
		return
	}
	// compress imageSt
	imageBytes = resizeImage(imageBytes)
	if len(imageBytes) == 0 {
		c.JSON(400, gin.H{
			"message": "Error compressing image",
		})
		return
	}
	stmt, err := db.Prepare("INSERT INTO images (id, image) VALUES (?, ?)")
	if err != nil {
		c.JSON(500, gin.H{
			"message": serverErrorMessage,
		})
		log.Println(err.Error())
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	var id = createRandomID()
	for checkImageID(id) {
		id = createRandomID()
	}
	_, err = stmt.Exec(id, imageBytes)
	if err != nil {
		c.JSON(500, gin.H{
			"message": serverErrorMessage,
		})
		log.Println(err.Error())
		return
	}
	c.JSON(201, gin.H{
		"message": "Image posted",
		"id":      id,
	})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func checkImageID(id string) bool {
	// Check if image id exists
	stmt, err := db.Prepare("SELECT id FROM images WHERE id = ?")
	if err != nil {
		return false
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var imageID string
	err = stmt.QueryRow(id).Scan(&imageID)
	return err == nil
}

func resizeImage(imageBytes []byte) []byte {
	// Resize imageData
	imageData, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return []byte{}
	}

	newImage := resize.Resize(1024, 1024, imageData, resize.Lanczos3)

	// Encode uses a Writer, use a Buffer if you need the raw []byte
	var bufferWriter = new(bytes.Buffer)
	err = jpeg.Encode(bufferWriter, newImage, nil)
	if err != nil {
		return []byte{}
	}
	return bufferWriter.Bytes()
}
