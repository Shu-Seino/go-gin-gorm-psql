package main

import (
	"database/sql"
	// "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

type Member struct {
    ID   int
    Name string
}


func membersHandler(c *gin.Context) {
    // データベース接続のセットアップ
    db, err := sql.Open("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
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

    // サーバーの起動
    r.Run(":8080")
}
