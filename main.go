package main

import (
	"mouse-macro/app"
	"mouse-macro/hook"
)

func main() {
	go hook.InitHook()
	app.InitApp()
}
