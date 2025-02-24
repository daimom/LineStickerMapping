package main

import (
	"database/sql"
	"fmt"

	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"log"
	// "strconv"
	// "strings"

	_ "github.com/mattn/go-sqlite3" // 匯入 SQLite 驅動
)

type ProductInfo struct {
	PackageID int `json:"packageId"`
	titles    []struct {
		title string `json:"zh-Hant"`
	} `json:"title"`
	Stickers []struct {
		ID int `json:"id"`
	} `json:"stickers"`
}

type Product struct {
	PackageID int
	title     string
	stickerId int
}

func searchAndParseFiles(db *sql.DB, root string) {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "productInfo.meta" {
			parseFile(db, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

func parseFile(db *sql.DB, filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", filePath, err)
		return
	}

	var info ProductInfo
	if err := json.Unmarshal(file, &info); err != nil {
		fmt.Println("Error parsing JSON:", filePath, err)
		return
	}
	var products []Product
	fmt.Printf("File: %s\nPackage ID: %d\n", filePath, info.PackageID)
	fmt.Printf("title: %s\n", info.titles[0].title)
	fmt.Print("Stickers IDs: ")
	for _, sticker := range info.Stickers {
		fmt.Printf("%d ", sticker.ID)
		products = append(products, Product{info.PackageID, info.titles[0].title, sticker.ID})
	}
	fmt.Println("\n--------------------")
	insertData(db, &products)
}

func main() {
	db, err := initDatabase("./stickersInfo.db")
	if err != nil {
		log.Fatalf("初始化資料庫失敗: %v", err)
	}
	defer db.Close()
	rootDir := "/Users/daimom/Library/Group Containers/VUTU7AKEUR.jp.naver.line.mac/Real/Library/Data/Sticker/33553" // 修改為你的目錄
	searchAndParseFiles(db, rootDir)
}

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
		packageId INTEGER NOT NULL,
		title TEXT NOT NULL,
		stickerSn INTEGER NOT NULL,
		stickerId INTEGER NOT NULL,
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
