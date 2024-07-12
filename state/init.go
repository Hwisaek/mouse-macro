package state

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/samber/lo"
)

var (
	MoveChecked = binding.BindBool(lo.ToPtr(false))
)
