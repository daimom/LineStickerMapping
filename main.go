package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // 匯入 SQLite 驅動

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	myWindow.SetContent(CreateUI(myApp))

	// 顯示視窗
	myWindow.ShowAndRun()

}

// 開啟圖片顯示的視窗
func openImageWindow(myApp fyne.App, imagePath string) {
	// 新視窗
	imageWindow := myApp.NewWindow("Image Viewer")

	// 顯示圖片
	img := canvas.NewImageFromFile(imagePath)
	img.FillMode = canvas.ImageFillOriginal

	// 顯示圖片
	imageWindow.SetContent(container.NewCenter(img))

	// 顯示圖片視窗
	imageWindow.Show()
}
