package main

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/mattn/go-sqlite3" // 匯入 SQLite 驅動
)

func main() {
	db, err := initDatabase("./stickersInfo.db")
	if err != nil {
		log.Fatalf("初始化資料庫失敗: %v", err)
	}
	defer db.Close()
	rootDir := "/Users/daimom/Library/Group Containers/VUTU7AKEUR.jp.naver.line.mac/Real/Library/Data/Sticker/33553" // 修改為你的目錄

	products, err := searchAndParseFiles(rootDir)
	if err != nil {
		fmt.Printf("Error search: %v\n", err)
	}
	insertData(db, products)
}

func insertData(db *sql.DB, products *[]Product) {

	// 使用 Transaction 提高效能
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO stickers (packageId, title, stickerSn, stickerId) 
				VALUES (?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	sn := 1
	for _, p := range *products {
		_, err = stmt.Exec(p.PackageID, p.title, sn, p.stickerId)
		if err != nil {
			log.Fatal(err)
		}
		sn++
	}

	// 提交事務
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Batch insert completed.")
}
