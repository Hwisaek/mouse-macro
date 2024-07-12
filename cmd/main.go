package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"mouse-macro/hook"
	"mouse-macro/layout"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Mouse macro")
	myWindow.Resize(fyne.Size{
		Width:  300,
		Height: 50,
	})
	myWindow.SetFixedSize(false)

	layout.Init(myWindow)
	hook.Init()

	myWindow.ShowAndRun()
}
