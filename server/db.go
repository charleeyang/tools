package main

import (
	"database/sql"
	_ "embed"
	"log"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

var db *sql.DB

func initDB(dsn string) {
	var err error
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	db.SetMaxOpenConns(1) // SQLite 单写入连接，避免 "database is locked"
	if _, err = db.Exec(schemaSQL); err != nil {
		log.Fatalf("建表失败: %v", err)
	}

	var n int
	if err = db.QueryRow("SELECT COUNT(*) FROM employee").Scan(&n); err != nil {
		log.Fatalf("检查种子数据失败: %v", err)
	}
	if n == 0 {
		log.Println("首次启动，写入种子数据 ...")
		seedAll()
		log.Println("种子数据写入完成")
	}
}
