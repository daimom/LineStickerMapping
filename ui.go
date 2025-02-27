package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 建立主 UI
func CreateUI(w fyne.App) fyne.CanvasObject {
	var imagesContainer *fyne.Container

	updateBtn := widget.NewButton("更新資料", func() {
		update()
	})

	searchBtn := widget.NewButton("搜尋", func() {
		fmt.Println("搜尋")
	})

	// 文字框
	textbox := widget.NewEntry()
	textbox.SetPlaceHolder("Enter Keyword...")

	// 頂部按鈕
	loadButton := widget.NewButton("載入圖片", func() {
		imageLists := read_packageID()
		imagesContainer.Objects = nil

		grid := container.NewGridWithColumns(3, LoadImages(*imageLists, w)...)
		imagesContainer.Objects = []fyne.CanvasObject{grid}
		// imagesContainer.Objects = LoadImages(*imagePaths, w)
		imagesContainer.Refresh()
	})

	// 頂部按鈕排版
	buttonContainer := container.NewGridWithColumns(4, updateBtn, loadButton, textbox, searchBtn)

	// 圖片列表
	imagesContainer = container.NewVBox()

	// 可滾動的圖片區域
	scroll := container.NewVScroll(imagesContainer)

	// 主界面佈局
	return container.NewBorder(buttonContainer, nil, nil, nil, scroll)
}

// 根據圖片路徑載入圖片
func LoadImages(imageLists []ImageInfo, parent fyne.App) []fyne.CanvasObject {
	var images []fyne.CanvasObject
	var path, title string
	for _, values := range imageLists {
		path = values.FolderPath + "tab_on@2x.png"
		title = values.Title

		img := canvas.NewImageFromFile(path)
		img.FillMode = canvas.ImageFillOriginal

		openButton := widget.NewButton("檢視", func() {
			ShowImageWindow(values.FolderPath, title, parent)
		})

		// 圖片和按鈕組合
		images = append(images, container.NewVBox(img, openButton))
	}
	return images
}

// 顯示新視窗
func ShowImageWindow(imagePath string, title string, parent fyne.App) {
	p := strings.Split(imagePath, "/")
	packageID := p[len(p)-2] //path 為 /abc/abc/123/ ，故-2
	stickerLists := read_stickerID(packageID)

	var images []fyne.CanvasObject
	w := parent.NewWindow("檢視圖片")
	for _, val := range *stickerLists {
		fulPath := imagePath + val + "_key@2x.png"
		img := canvas.NewImageFromFile(fulPath)
		img.FillMode = canvas.ImageFillOriginal

		aliasButton := widget.NewButton("alias", func() {
			ShowAliasWindow(fulPath, parent)
		})
		// 圖片和按鈕組合
		images = append(images, container.NewVBox(img, aliasButton))

	}
	//img := canvas.NewImageFromFile("")
	grid := container.NewGridWithColumns(3, images...)

	w.SetContent(container.NewVScroll(grid))
	w.Resize(fyne.NewSize(400, 400))
	w.Show()
}

// 顯示Alias視窗
func ShowAliasWindow(filePath string, parent fyne.App) {
	p := strings.Split(filePath, "/")
	fileName := p[len(p)-1] //path 為 /abc/abc/123/ ，故-2
	stickerId := strings.ReplaceAll(fileName, "_key@2x.png", "")

	w := parent.NewWindow("新增別名")

	img := canvas.NewImageFromFile(filePath)
	img.FillMode = canvas.ImageFillOriginal

	// Get sticker alias
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter values, separated by commas.")

	insertButton := widget.NewButton("新增", func() {
		fmt.Print(stickerId)
		//insert sticker alias
		//ShowImageWindow(values.FolderPath, title, parent)
	})

	top := container.NewCenter(img)
	// center := container.NewVBox(input)
	bottom := container.NewVBox(insertButton)
	w.SetContent(container.NewBorder(top, bottom, nil, nil, input))
	w.Resize(fyne.NewSize(400, 400))
	w.Show()
}
