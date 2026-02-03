package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
	ServerPort string `json:"server_port"`
}

var db *sql.DB
var appConfig Config

// 【修改：增加读取配置文件函数】
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("无法打开配置文件 config.json:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appConfig)
	if err != nil {
		log.Fatal("配置文件格式错误:", err)
	}
}

func initDB() {
	// DSN格式: 用户名:密码@tcp(地址:3306)/数据库名
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		appConfig.DBUser,
		appConfig.DBPassword,
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBName,
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 测试数据库连接
	if err = db.Ping(); err != nil {
		log.Fatal("数据库连接失败，请检查配置:", err)
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
	loadConfig()
	initDB()
	defer db.Close()

	// 1. 路由：访问 /api/visit 时调用接口
	http.HandleFunc("/api/visit", visitHandler)

	// 2. 路由：访问根目录 / 时，自动加载并展示 index.html
	// http.FileServer 会自动寻找目录下的 index.html
	http.Handle("/", http.FileServer(http.Dir("./")))

	fmt.Printf("🚀 服务已启动: http://localhost%s\n", appConfig.ServerPort)

	// 启动监听
	log.Fatal(http.ListenAndServe(appConfig.ServerPort, nil))
}
