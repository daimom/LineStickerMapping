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
	rows, err := db.Query("SELECT distinct folderpath,title FROM stickers ")
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
func readKeyword(keyword string) *[]ImageInfo {
	keyword = "%" + keyword + "%"
	// 查詢 stickers 資料表中的 folderpath 欄位
	rows, err := db.Query(`select distinct stickers.folderpath,stickers.title 
from  stickers  inner join alias on stickers.stickerId = alias.stickerId 
where alias.alias like ?`, keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var imageLists []ImageInfo
	// 遍歷查詢結果
	for rows.Next() {
		var filePath, title string
		if err := rows.Scan(&filePath, &title); err != nil {
			log.Fatal(err)
		}
		imageLists = append(imageLists, ImageInfo{FolderPath: filePath, Title: title})

	}

	// 檢查查詢過程中的錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &imageLists
}

func read_stickerID(packageId string) *[]string {

	// 查詢 stickers 資料表中的 folderpath 欄位
	rows, err := db.Query(`SELECT stickerId FROM stickers where packageId=?`, packageId)
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
func readAlias(stickerId string) *[]string {

	// 查詢 stickers 資料表中的 folderpath 欄位
	rows, err := db.Query("SELECT alias FROM alias where stickerId=?", stickerId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var aliasLists []string
	// 遍歷查詢結果
	for rows.Next() {
		var alias string
		if err := rows.Scan(&alias); err != nil {
			log.Fatal(err)
		}
		aliasLists = append(aliasLists, alias)

	}

	// 檢查查詢過程中的錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &aliasLists
}
func deleteAlias(stickerId string) error {
	// 使用 Transaction 提高效能
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	stmt, err := tx.Prepare(`delete from alias where stickerId=?`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(stickerId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute statement for stickerId %s: %w", stickerId, err)
	}

	// 提交事務
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
func insertAlias(stickers *[]Sticker) error {

	// 使用 Transaction 提高效能
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO alias (stickerId,alias)
				VALUES (? , ?)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var lastStickerId string
	for _, p := range *stickers {
		lastStickerId = p.StickerId
		_, err = stmt.Exec(p.StickerId, p.Alias)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement for stickerId %s: %w", p.StickerId, err)
		}
	}

	// 提交事務
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Alias %s insert completed.", lastStickerId)
	return nil
}
