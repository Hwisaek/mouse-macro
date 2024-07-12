package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"log"
	"math/rand"
	"mouse-macro/windows/mouse"
	"time"
)

func main() {
	myApp := app.New()

	myWindow := myApp.NewWindow("Mouse macro")
	myWindow.Resize(fyne.Size{
		Width:  300,
		Height: 50,
	})
	myWindow.SetFixedSize(false)

	checkedData := false
	moveChecked := binding.BindBool(&checkedData)
	check := widget.NewCheckWithData("move", moveChecked)
	check.OnChanged = func() func(checked bool) {
		var ticker *time.Ticker
		var done chan bool

		return func(checked bool) {
			if checked {
				// 체크박스가 체크되면 타이머 시작
				ticker = time.NewTicker(1 * time.Second)
				done = make(chan bool)
				go func() {
					for {
						select {
						case <-done:
							return
						case <-ticker.C:
							x, y, err := mouse.Location()
							if err != nil {
								log.Println(err)
								return
							}

							newX := x + (rand.Int()%10 - 5)
							newY := y + (rand.Int()%10 - 5)
							robotgo.Move(newX, newY)
						}
					}
				}()
			} else {
				// 체크박스가 해제되면 타이머 중지
				if ticker != nil {
					ticker.Stop()
					done <- true
				}
			}
		}
	}()

	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		checked, err := moveChecked.Get()
		if err != nil {
			log.Println(err)
			return
		}

		moveChecked.Set(!checked)
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()

	myWindow.SetContent(container.NewVBox(
		check,
	))

	myWindow.ShowAndRun()
}
