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
	checkDuplicate(products)
	insertData(products)
}

func insertData(products *[]Product) {

	// 使用 Transaction 提高效能
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// duplicate
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
func checkDuplicate(products *[]Product) {

	// 宣告 空結構體
	packageIDSet := make(map[int]struct{})
	for _, p := range *products {
		packageIDSet[p.PackageID] = struct{}{}
	}

	// 將 map 轉成 slice
	var packageIDs []int
	for pkg := range packageIDSet {
		packageIDs = append(packageIDs, pkg)
	}

	if len(packageIDs) == 0 {
		return
	}

	// 構造 SQL 查詢 (IN 條件式)
	query := "SELECT DISTINCT packageId FROM stickers WHERE packageId IN ("
	args := make([]interface{}, len(packageIDs))
	for i, id := range packageIDs {
		query += "?"
		if i < len(packageIDs)-1 {
			query += ","
		}
		args[i] = id
	}
	query += ")"

	// 執行查詢
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("查詢 package_id 發生錯誤: %v", err)
	}
	defer rows.Close()

	// 把 DB 存在的 packageId 存起來
	existsMap := make(map[int]bool)
	for rows.Next() {
		var pkgID int
		if err := rows.Scan(&pkgID); err != nil {
			log.Fatalf("掃描 package_id 發生錯誤: %v", err)
		}
		existsMap[pkgID] = true
	}

	// 過濾 products：只留下那些 package_id 沒在 existsMap 裡的
	filtered := (*products)[:0] // in-place 過濾，建立空的slice(可避免重新分配記憶體)
	for _, p := range *products {
		if !existsMap[p.PackageID] {
			filtered = append(filtered, p)
		} else {
			fmt.Printf("刪除 PackageID: %d, StickerID: %d\n", p.PackageID, p.stickerId)
		}
	}

	*products = filtered
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
