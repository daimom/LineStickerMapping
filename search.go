package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Product struct {
	PackageID int
	title     string
	stickerId int
}

type ProductInfo struct {
	PackageID int `json:"packageId"`
	Titles    struct {
		ZhHant string `json:"zh-Hant"`
	} `json:"title"`
	Stickers []struct {
		ID int `json:"id"`
	} `json:"stickers"`
}

func searchAndParseFiles(root string) (*[]Product, error) {
	var products []Product // 定義這個變數來存儲結果
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "productInfo.meta" {
			parsedProducts, err := parseFile(path)
			if err != nil {
				return fmt.Errorf("Error search %v", err)
			}
			products = append(products, *parsedProducts...) // 追加到產品列表

			return nil
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
	return &products, err
}

func parseFile(filePath string) (*[]Product, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file:%s: %v", filePath, err)
	}

	var info ProductInfo
	if err := json.Unmarshal(file, &info); err != nil {
		return nil, fmt.Errorf("Error parsing JSON:%s: %v", filePath, err)
	}
	var products []Product
	fmt.Printf("File: %s\nPackage ID: %d\n", filePath, info.PackageID)
	fmt.Printf("title: %s\n", info.Titles.ZhHant)
	fmt.Print("Stickers IDs: ")
	for _, sticker := range info.Stickers {
		fmt.Printf("%d ", sticker.ID)
		products = append(products, Product{info.PackageID, info.Titles.ZhHant, sticker.ID})
	}
	fmt.Println("\n--------------------")
	return &products, nil
}
