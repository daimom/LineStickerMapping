package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // 匯入 SQLite 驅動

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var BasePath = "/Users/daimom/Library/Group Containers/VUTU7AKEUR.jp.naver.line.mac/Real/Library/Data/Sticker/"
var db *sql.DB

func main() {
	var err error
	db, err = initDatabase("./stickersInfo.db")
	if err != nil {
		log.Fatalf("初始化資料庫失敗: %v", err)
	}
	defer db.Close()

	myApp := app.New()
	myWindow := myApp.NewWindow("Line Sticker Alias")

	// 設定視窗的初始大小
	myWindow.Resize(fyne.NewSize(800, 600))

	// 設定視窗的內容
	myWindow.SetContent(CreateUI(myWindow))

	// 顯示視窗
	myWindow.ShowAndRun()

}
