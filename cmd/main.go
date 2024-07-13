package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
	"mouse-macro/hook"
	"mouse-macro/layout"
	"os"
)

func main() {
	myApp := app.New()

	icon, err := os.ReadFile("icon.png")
	if err != nil {
		log.Println("failed to load icon :", err)
	}

	myApp.SetIcon(fyne.NewStaticResource("icon.png", icon))

	myWindow := myApp.NewWindow("mouse-macro")
	myWindow.Resize(fyne.Size{Width: 300, Height: 550})
	myWindow.SetFixedSize(false)

	layout.Init(myWindow)
	hook.Init()

	myWindow.ShowAndRun()
}
