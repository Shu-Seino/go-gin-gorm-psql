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
	FavoriteFood string
}

var db *gorm.DB

func init() {
	godotenv.Load(".devcontainer/.env")
	POSTGRES_USER,POSTGRES_PASSWORD,POSTGRES_DB,POSTGRES_HOSTNAME := os.Getenv("POSTGRES_USER"),os.Getenv("POSTGRES_PASSWORD"),os.Getenv("POSTGRES_DB"),os.Getenv("POSTGRES_HOSTNAME_")

	conn, err := gorm.Open("postgres", "host="+POSTGRES_HOSTNAME+" dbname="+POSTGRES_DB+" password="+POSTGRES_PASSWORD+" user="+POSTGRES_USER+" sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db = conn
	db.AutoMigrate(&Member{})
}

func insertHandler(c *gin.Context) {
	name := c.PostForm("name")
	favoriteFood := c.PostForm("favoriteFood")

	member := Member{Name: name, FavoriteFood: favoriteFood}
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

	if err := db.Where("id = ?", id).Delete(&Member{}).Error; err != nil {
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
	viewModel :=  ListViewModel{
		MemberList: members,
	}
	c.HTML(200, "members.html",viewModel)
}
type ListViewModel struct {
	MemberList []Member ; 
}
func memberHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		c.String(400, "Invalid ID")
		return
	}

	var member Member
	if err := db.First(&member, id).Error; err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}

	c.HTML(200, "member.html", gin.H{
		"Member": member,
	})
}

func main() {
	r := gin.Default()

	r.Static("/static","./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/members", membersHandler)
	r.POST("/members", insertHandler)
	r.POST("/delete/members/:id", deleteHandler) 
	r.GET("/members/:id", memberHandler)

	r.Run(":8080")
}