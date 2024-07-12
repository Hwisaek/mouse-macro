package monitor

import (
	"awesomeProject/windows"
	"fmt"
	"sort"
	"syscall"
	"unsafe"
)

var (
	procEnumDisplayMonitors = windows.User32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfo      = windows.User32.NewProc("GetMonitorInfoW")
)

type monitorInfo struct {
	CbSize    uint32
	RcMonitor Monitor
	RcWork    Monitor
	DwFlags   uint32
}

type Monitor struct {
	Left, Top, Right, Bottom int32
}

var MonitorList []Monitor

func init() {
	GetMonitorsInfo()
}
func GetMonitorsInfo() (err error) {
	callback := syscall.NewCallback(func(hMonitor, hdcMonitor, lprcMonitor, dwData uintptr) uintptr {
		monitorInfo := monitorInfo{CbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
		ret, _, _ := procGetMonitorInfo.Call(hMonitor, uintptr(unsafe.Pointer(&monitorInfo)))
		if ret != 0 {
			MonitorList = append(MonitorList, monitorInfo.RcMonitor)
		}

		//order by monitorList
		sort.Slice(MonitorList, func(i, j int) bool {
			return MonitorList[i].Left < MonitorList[j].Left
		})

		return 1 // Continue enumeration
	})

	ret, _, _ := procEnumDisplayMonitors.Call(0, 0, callback, 0)
	if ret == 0 {
		return fmt.Errorf("GetMonitorsInfo failed")
	}
	return nil
}
