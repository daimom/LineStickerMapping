package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// 建立主 UI
func CreateUI(w fyne.App) fyne.CanvasObject {
	var imagesContainer *fyne.Container

	updateBtn := widget.NewButton("更新資料", func() {
		update()
	})

	// 文字框
	textbox := widget.NewEntry()
	textbox.SetPlaceHolder("Enter Keyword...")

	searchBtn := widget.NewButton("搜尋", func() {
		imageLists := readKeyword(textbox.Text)

		fmt.Println("搜尋：", textbox.Text)
	})

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
			time.AfterFunc(200*time.Millisecond, func() {
				ShowImageWindow(values.FolderPath, title, parent)
			})

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
			time.AfterFunc(200*time.Millisecond, func() {
				ShowAliasWindow(fulPath, parent)
			})
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

type Sticker struct {
	StickerId string
	Alias     string
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
	aliasList := readAlias(stickerId)
	if len(*aliasList) == 0 {
		input.SetPlaceHolder("Enter values, separated by commas.")
	} else {
		input.Text = (strings.Join(*aliasList, ","))
	}

	insertButton := widget.NewButton("新增", func() {
		fmt.Print(stickerId)
		//insert sticker alias
		inputText := input.Text
		stickers := parseInput(inputText, stickerId)
		// // 顯示解析後的結構
		// for _, sticker := range stickers {
		// 	fmt.Printf("StickerId: %s, Alias: %s\n", sticker.StickerId, sticker.Alias)
		// }
		err := deleteAlias(stickerId)
		if err != nil {
			dialog.NewError(err, w).Show()
			return // 錯誤發生後，直接返回，停止繼續執行下面的程式碼
		}
		err2 := insertAlias(&stickers)
		if err2 != nil {
			// 若出現錯誤，顯示錯誤訊息並不關閉視窗
			dialog.NewError(err2, w).Show()
		} else {
			// 如果成功，顯示提示框後關閉視窗
			dialog.NewInformation("提示", "新增完成！", w).Show()

			go func() {
				<-time.After(2 * time.Second) // 延遲 2 秒後關閉視窗
				w.Close()                     // 關閉視窗
			}()
		}
	})

	top := container.NewCenter(img)
	// center := container.NewVBox(input)
	bottom := container.NewVBox(insertButton)
	w.SetContent(container.NewBorder(top, bottom, nil, nil, input))
	w.Resize(fyne.NewSize(400, 400))
	w.Show()

}

// parseInput 解析輸入的字串，並返回 Sticker struct 陣列
func parseInput(input string, strickerId string) []Sticker {
	var stickers []Sticker
	// 拆分輸入的字串，假設每個 sticker 的 stickerId 和 alias 由 ',' 隔開
	entries := strings.Split(input, ",")

	// 檢查每對 stickerId 和 alias
	for _, p := range entries {
		sticker := Sticker{
			StickerId: strickerId, // 去除空格
			Alias:     strings.TrimSpace(p),
		}
		stickers = append(stickers, sticker)
	}
	return stickers
}
