package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	"mouse-macro/process"
	"mouse-macro/process/memory"
	"unsafe"
)

func main() {
	pid, found := process.GetPid("winmine_(한글).exe")
	if !found {
		log.Fatalln("process not found")
	}

	handle, err := process.Open(pid)
	if err != nil {
		log.Fatalln("process handle not found")
	}
	defer process.Close()

	process.Handle = handle

	var startAddress uintptr = 0x0000000000000000
	var endAddress uintptr = 0x7FFFFFFFFFFFFFFF
	_ = endAddress

	type MemoryInfo struct {
		Addr uintptr
		Data uint32
	}
	validMemoryList := make([]MemoryInfo, 0)

	// 메모리 영역 정보를 저장할 구조체
	var mbi windows.MemoryBasicInformation

	var targetAddress uintptr = 0x0100579C

	// 메모리 주소 범위 순회
	for addr := startAddress; ; {
		// 메모리 영역 정보 조회
		err := windows.VirtualQueryEx(handle, addr, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			break // 더 이상 조회할 메모리 영역이 없으면 종료
		}

		if mbi.State == windows.MEM_COMMIT {
			for addr := mbi.BaseAddress; addr < mbi.BaseAddress+mbi.RegionSize; {
				data, bytesRead := memory.ReadMemory(addr)
				if bytesRead == 0 {
					break
				}

				// 특정 값 찾기
				if addr == targetAddress {
					fmt.Println("idx: ", len(validMemoryList))
					fmt.Printf("found value at: 0x%X, data: %v\n", addr, data)
				}

				m := MemoryInfo{
					Addr: addr,
					Data: data,
				}

				validMemoryList = append(validMemoryList, m)
				addr += bytesRead
			}
		}

		if startAddress == 0 {
			startAddress = mbi.BaseAddress
		}

		addr = mbi.BaseAddress + mbi.RegionSize
		endAddress = addr
	}

	// 메모리 순회
	for _, mbi := range validMemoryList {
		if mbi.Addr == targetAddress {
			fmt.Printf("found value at: 0x%X, data: %v\n", mbi.Addr, mbi.Data)
			break
		}

	}

	data, _ := memory.ReadMemory(targetAddress)
	fmt.Println(data)

	if !process.HasAdminPermission() {
		fmt.Println("This program needs to be run as administrator.")
		return
	}

	memory.WriteMemory(targetAddress, 42)
}
