package hook

import (
	"math/rand"
	"mouse-macro/status"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func InitHook() {
	hook.Register(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
		DeactivateMacro()
	})

	// 이 부분에서 프로그램을 종료하지 않고 대기합니다.
	// ESC 키를 누를 때까지 계속해서 이벤트를 감지합니다.
	s := hook.Start()
	<-hook.Process(s)
}

func ActivateMacro() {
	status.MacroFlag.Set(true)

	for {
		flag, _ := status.MacroFlag.Get()
		if !flag {
			break
		}

		x, y := robotgo.Location()
		dx, dy := rand.Intn(20)-10, rand.Intn(20)-10

		robotgo.Move(x+dx, y+dy)

		time.Sleep(time.Second)
	}
}

func DeactivateMacro() {
	status.MacroFlag.Set(false)
}
