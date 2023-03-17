package app

import (
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

func InitApp() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewButton("Activate", clickBtn_Macro()))
	w.ShowAndRun()
}

func clickBtn_Macro() func() {
	flag := false
	return func() {
		flag = !flag
		go func() {
			for flag {
				robotgo.Move(100, 100)
				time.Sleep(time.Second)
				robotgo.Move(200, 200)
				time.Sleep(time.Second)
			}
		}()
	}
}
