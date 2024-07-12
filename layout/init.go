package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"log"
	"math/rand/v2"
	"mouse-macro/state"
	"mouse-macro/windows/mouse"
	"time"
)

func Init(window fyne.Window) {
	check := widget.NewCheckWithData("move", state.MoveChecked)
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

	window.SetContent(container.NewVBox(
		check,
	))

}
