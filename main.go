package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
	"log"
	"strconv"
)

type Member struct {
	gorm.Model
	Name string
}

var db *gorm.DB

func init() {
	godotenv.Load(".env")
	POSTGRES_USER,POSTGRES_PASSWORD,POSTGRES_DB := os.Getenv("POSTGRES_USER"),os.Getenv("POSTGRES_PASSWORD"),os.Getenv("POSTGRES_DB")
	conn, err := gorm.Open("postgres", "user="+POSTGRES_USER+" dbname="+POSTGRES_DB+" password="+POSTGRES_PASSWORD+" sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db = conn
	db.AutoMigrate(&Member{})
}

func insertHandler(c *gin.Context) {
	name := c.PostForm("name")

	member := Member{Name: name}
	if err := db.Create(&member).Error; err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}

	c.Redirect(302, "/members")
}

func deleteHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		c.String(400, "Invalid ID")
		return
	}

	if err := db.Delete(&Member{}, id).Error; err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}

	c.Redirect(302, "/members")
}

func membersHandler(c *gin.Context) {
	var members []Member
	if err := db.Find(&members).Error; err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}

	c.HTML(200, "members.html", gin.H{
		"MemberList": members,
	})
}

func main() {
	r := gin.Default()

	r.Static("/static","./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/members", membersHandler)
	r.POST("/members", insertHandler)
	r.DELETE("/members/:id", deleteHandler)

	r.Run(":8080")
}