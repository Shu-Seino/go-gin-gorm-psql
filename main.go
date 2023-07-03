package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"os"
	// "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

type Member struct {
	ID   int
	Name string
}

func insertMember(db *sql.DB, id int, name string) error {
	_, err := db.Exec("INSERT INTO member (id, name) VALUES ($1, $2)", id, name)
	if err != nil {
		return err
	}
	return nil
}

func insertHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	name := c.PostForm("name")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		c.String(400, "Invalid ID")
		return
	}

	// データベース接続のセットアップ
	godotenv.Load(".env")
	POSTGRES_USER,POSTGRES_PASSWORD,POSTGRES_DB := os.Getenv("POSTGRES_USER"),os.Getenv("POSTGRES_PASSWORD"),os.Getenv("POSTGRES_DB")
	db, err := sql.Open("postgres", "user="+POSTGRES_USER+" dbname="+POSTGRES_DB+" password="+POSTGRES_PASSWORD+" sslmode=disable")
	if err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}
	defer db.Close()

	// メンバーの追加
	err = insertMember(db, id, name)
	if err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}

	// メンバーリストの表示ページにリダイレクト
	c.Redirect(302, "/members")
}

func membersHandler(c *gin.Context) {
	// データベース接続のセットアップ
    godotenv.Load(".env")
	POSTGRES_USER,POSTGRES_PASSWORD,POSTGRES_DB := os.Getenv("POSTGRES_USER"),os.Getenv("POSTGRES_PASSWORD"),os.Getenv("POSTGRES_DB")
	db, err := sql.Open("postgres", "user="+POSTGRES_USER+" dbname="+POSTGRES_DB+" password="+POSTGRES_PASSWORD+" sslmode=disable")
	if err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}
	defer db.Close()

	// メンバーリストの取得
	rows, err := db.Query("SELECT * FROM member")
	if err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}
	defer rows.Close()

	memberList := []Member{}

	for rows.Next() {
		var member Member
		err := rows.Scan(&member.ID, &member.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		memberList = append(memberList, member)
	}

	// HTMLテンプレートにデータを渡してレンダリング
	c.HTML(200, "members.html", gin.H{
		"MemberList": memberList,
	})
}

func main() {
	r := gin.Default()

	// 静的ファイルの配信
	r.Static("/static", "./static")

	// HTMLテンプレートのロード
	r.LoadHTMLGlob("templates/*")

	// メンバーリストのルーティングとハンドラの設定
	r.GET("/members", membersHandler)
	r.POST("/members", insertHandler)

	// サーバーの起動
	r.Run(":8080")
}
