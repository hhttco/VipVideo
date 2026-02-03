package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	// DSN格式: 用户名:密码@tcp(地址:3306)/数据库名
	dsn := "video:vip#video123!@tcp(127.0.0.1:3306)/video?charset=utf8mb4&parseTime=True"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
}

// 记录并获取访问数据
func visitHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	// 插入数据库
	db.Exec("INSERT INTO visitor_logs (ip_address, user_agent) VALUES (?, ?)", ip, r.UserAgent())

	// 获取总数
	var count int
	db.QueryRow("SELECT COUNT(*) FROM visitor_logs").Scan(&count)

	// 返回 JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_visits": count,
		"ip":           ip,
	})
}

func main() {
	initDB()
	defer db.Close()

	// 1. 路由：访问 /api/visit 时调用接口
	http.HandleFunc("/api/visit", visitHandler)

	// 2. 路由：访问根目录 / 时，自动加载并展示 index.html
	// http.FileServer 会自动寻找目录下的 index.html
	http.Handle("/", http.FileServer(http.Dir("./")))

	fmt.Println("🚀 服务启动成功！")
	fmt.Println("🔗 请访问: http://localhost:8080")

	// 启动监听
	log.Fatal(http.ListenAndServe(":8080", nil))
}
