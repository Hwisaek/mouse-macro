package memory

import (
	"golang.org/x/sys/windows"
	"mouse-macro/process"
	"unsafe"
)

func ReadMemory(address uintptr) (data uint32, bytesRead uintptr) {
	err := windows.ReadProcessMemory(process.Handle, address, (*byte)(unsafe.Pointer(&data)), unsafe.Sizeof(data), &bytesRead)
	if err != nil {
		return
	}

	return
}
