package hook

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

func InitHook() {
	hook.Register(hook.KeyDown, []string{"w"}, func(hook.Event) {
		fmt.Println("w")
	})

	go func() {
		<-hook.Process(hook.Start())
	}()
}
