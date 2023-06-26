package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		fmt.Println("エラー")
	}
	defer db.Close()
	rows, err := db.Query("select * from member")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	var id int
	var name string

	for rows.Next() {
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": name,
		})
	})
	r.Run()

}
