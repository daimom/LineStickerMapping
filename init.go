package main

import (
	"database/sql"
	"fmt"
)

func initDatabase(dbPath string) (*sql.DB, error) {
	// 開啟或創建 SQLite 資料庫
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("無法開啟資料庫: %v", err)
	}
	// 創建表格
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS stickers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		folderpath TEXT NOU NULL,
		packageId INTEGER NOT NULL,
		title TEXT NOT NULL,
		stickerSn INTEGER NOT NULL,
		stickerId INTEGER NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("創建stickers表格失敗: %v", err)
	}

	createTableSQL = `
	CREATE TABLE IF NOT EXISTS alias (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stickerId INTEGER NOT NULL,
		alias TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("創建alias表格失敗: %v", err)
	}
	fmt.Println("表格 'stickers' 已成功創建或已存在。")
	return db, nil
}
