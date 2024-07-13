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
	if !process.HasAdminPermission() {
		fmt.Println("This program needs to be run as administrator.")
		return
	}

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

	var targetAddress uintptr = 0x0100579C // 지뢰찾기 시간 메모리값

	// 메모리 주소 범위 순회
	for addr := startAddress; ; addr = mbi.BaseAddress + mbi.RegionSize {
		// 메모리 영역 정보 조회
		err := windows.VirtualQueryEx(handle, addr, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			break // 더 이상 조회할 메모리 영역이 없으면 종료
		}

		if startAddress == 0 {
			startAddress = mbi.BaseAddress
		}
		endAddress = mbi.BaseAddress

		if addr != mbi.BaseAddress {
			fmt.Println(addr)
		}

		switch mbi.Type {
		case 0x1000000: // MEM_IMAGE: 영역 내의 메모리 페이지가 이미지 섹션의 보기에 매핑됨을 나타냅니다.
		case 0x40000: // MEM_MAPPED: 영역 내의 메모리 페이지가 섹션 보기에 매핑됨을 나타냅니다.
		case 0x20000: // MEM_PRIVATE: 지역 내의 메모리 페이지가 프라이빗(즉, 다른 프로세스에서 공유되지 않음)임을 나타냅니다.
		default:
			continue
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

	memory.WriteMemory(targetAddress, 42)
}
