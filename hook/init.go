package hook

import (
	"fmt"
	hook "github.com/robotn/gohook"
	"log"
	"mouse-macro/state"
)

func Init() {
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		checked, err := state.MoveChecked.Get()
		if err != nil {
			log.Println(err)
			return
		}

		state.MoveChecked.Set(!checked)
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()
}
