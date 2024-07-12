package mouse

import (
	"fmt"
	"mouse-macro/windows"
	"mouse-macro/windows/monitor"
	"unsafe"
)

var (
	procSetCursorPos = windows.User32.NewProc("SetCursorPos")
	procGetCursorPos = windows.User32.NewProc("GetCursorPos")
)

func Move(x, y int32) error {
	ret, _, _ := procSetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return fmt.Errorf("SetCursorPos failed")
	}
	return nil
}

func Location() (x, y int, err error) {
	var pt point
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		err = fmt.Errorf("location failed")
		return
	}

	for _, m := range monitor.MonitorList {
		if m.Left > 0 {
			break
		}
		pt.X -= m.Left
	}

	return int(pt.X), int(pt.Y), nil
}

type point struct {
	X, Y int32
}

// 현재 마우스 위치가 어느 모니터에 있는지 확인하는 함수
func GetCurrentMonitor() (idx int, m monitor.Monitor, err error) {
	x, y, err := Location()
	if err != nil {
		return -1, m, err
	}

	for i, m := range monitor.MonitorList {
		if x >= int(m.Left) && x < int(m.Right) && y >= int(m.Top) && y < int(m.Bottom) {
			return i, m, nil
		}
	}

	err = fmt.Errorf("현재 마우스 위치에 해당하는 모니터를 찾을 수 없습니다")
	return
}
