package main

import (
	"mouse-macro/app"
	"mouse-macro/hook"
)

func main() {
	hook.InitHook()
	app.InitApp()
}
