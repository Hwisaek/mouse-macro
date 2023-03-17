package app

import (
	"mouse-macro/hook"
	"mouse-macro/status"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func InitApp() {
	a := app.New()

	w := a.NewWindow("Mouse Macro")
	w.Resize(fyne.NewSize(200, 100))

	chk_Macro := widget.NewCheckWithData("Macro", status.MacroFlag)
	chk_Macro.OnChanged = clickChk

	w.SetContent(chk_Macro)

	w.ShowAndRun()
}

func clickChk(b bool) {
	if b {
		hook.ActivateMacro()
	} else {
		hook.DeactivateMacro()
	}
}
