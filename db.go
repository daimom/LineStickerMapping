package main

import (
	"fmt"
	"log"
)

type ImageInfo struct {
	FolderPath string
	Title      string
}

func read_packageID() *[]ImageInfo {

	// 查詢 stickers 資料表中的 folderpath 欄位
	rows, err := db.Query("SELECT distinct folderpath,title FROM stickers limit 20")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var imageLists []ImageInfo
	// 遍歷查詢結果
	for rows.Next() {
		var folderPath, title string
		if err := rows.Scan(&folderPath, &title); err != nil {
			log.Fatal(err)
		}
		// 將 "/abc/def" 加在 packageID 前面
		fullPath := folderPath
		imageLists = append(imageLists, ImageInfo{FolderPath: fullPath, Title: title})

	}

	// 檢查查詢過程中的錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &imageLists
}

func read_stickerID(packageId string) *[]string {

	// 查詢 stickers 資料表中的 folderpath 欄位
	rows, err := db.Query("SELECT stickerId FROM stickers where packageId=?", packageId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var imageLists []string
	// 遍歷查詢結果
	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			log.Fatal(err)
		}
		imageLists = append(imageLists, filePath)

	}

	// 檢查查詢過程中的錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &imageLists
}

func update() {
	products, err := searchAndParseFiles(BasePath)
	if err != nil {
		fmt.Printf("Error search: %v\n", err)
	}
	insertData(products)
}

func insertData(products *[]Product) {

	// 使用 Transaction 提高效能
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO stickers (packageId,folderPath, title, stickerSn, stickerId) 
				VALUES (?,?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, p := range *products {
		_, err = stmt.Exec(p.PackageID, p.folderpath, p.title, p.stickerSn, p.stickerId)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 提交事務
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Batch insert completed.")
}

// func insertAlias(stickers string) {

// 	// 使用 Transaction 提高效能
// 	tx, err := db.Begin()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	stmt, err := tx.Prepare(`INSERT INTO alias (stickerId,alias)
// 				VALUES (? , ?)`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()

// 	for _, p := range *products {
// 		_, err = stmt.Exec(p.PackageID, p.folderpath, p.title, p.stickerSn, p.stickerId)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	// 提交事務
// 	err = tx.Commit()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Batch insert completed.")
// }
