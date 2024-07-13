package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/0xrawsec/golang-win32/win32"
	"github.com/0xrawsec/golang-win32/win32/kernel32"
	"golang.org/x/sys/windows"
)

var Handle windows.Handle
var procReadProcessMemory *windows.Proc
var baseAddress int64

func memoryReadInit(pid uint32) (int64, bool) {
	Handle, _ = windows.OpenProcess(kernel32.PROCESS_ALL_ACCESS, false, pid)
	procReadProcessMemory = windows.MustLoadDLL("kernel32.dll").MustFindProc("ReadProcessMemory")

	win32handle, _ := kernel32.OpenProcess(0x0010|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, win32.BOOL(0), win32.DWORD(pid))
	moduleHandles, _ := kernel32.EnumProcessModules(win32handle)
	for _, moduleHandle := range moduleHandles {
		s, _ := kernel32.GetModuleFilenameExW(win32handle, moduleHandle)
		targetModuleFilename := "UE4Game-Win64-Shipping.exe"
		if filepath.Base(s) == targetModuleFilename {
			info, _ := kernel32.GetModuleInformation(win32handle, moduleHandle)
			baseAddress = int64(info.LpBaseOfDll)
			return baseAddress, true
		}
	}
	return 0, false
}

func memoryReadClose() {
	windows.CloseHandle(Handle)
}

func readMemoryAt(address int64) float32 {
	var (
		data   [4]byte
		length uint32
	)

	procReadProcessMemory.Call(
		uintptr(Handle),
		uintptr(address),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&length)),
	)

	bits := binary.LittleEndian.Uint32(data[:])
	float := math.Float32frombits(bits)
	return float
}

func readMemory(address uintptr) (data uint32, bytesRead uintptr) {
	// 메모리 읽기
	err := windows.ReadProcessMemory(Handle, address, (*byte)(unsafe.Pointer(&data)), unsafe.Sizeof(data), &bytesRead)
	if err != nil {
		return
	}

	return
}

func readMemoryAtByte8(address int64) uint64 {
	var (
		data   [8]byte
		length uint32
	)

	procReadProcessMemory.Call(
		uintptr(Handle),
		uintptr(address),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&length)),
	)

	byte8 := binary.LittleEndian.Uint64(data[:])
	return byte8
}

type staticPointer struct {
	baseOffset int64
	offsets    []string
}

func GetAddresses() (int64, int64) {
	xPositionPointer := staticPointer{0x2518790, []string{"2E4", "10", "8", "8", "8", "78", "5E0"}}
	zPositionPointer := staticPointer{0x2518790, []string{"2E8", "10", "8", "8", "8", "78", "5E0"}}
	xPositionAddress := calculateAddress(xPositionPointer)
	zPositionAddress := calculateAddress(zPositionPointer)
	xPositionAddressInt, _ := strconv.ParseInt(xPositionAddress, 16, 0)
	zPositionAddressInt, _ := strconv.ParseInt(zPositionAddress, 16, 0)
	return xPositionAddressInt, zPositionAddressInt
}

func calculateAddress(pointer staticPointer) string {
	startingPointer := baseAddress + pointer.baseOffset
	startingAddress := readMemoryAtByte8(startingPointer)
	var value string = strconv.FormatInt(int64(startingAddress), 16)

	for i := len(pointer.offsets) - 1; i >= 0; i-- {
		offset := pointer.offsets[i]
		addressPointer := sumHex(value, offset)

		if i > 0 {
			addressInt, _ := strconv.ParseInt(addressPointer, 16, 64)
			nextAddressDecimal := readMemoryAtByte8(addressInt)
			value = strconv.FormatInt(int64(nextAddressDecimal), 16)
		} else {
			value = addressPointer
		}
	}
	return value
}

func sumHex(aHex string, bHex string) string {
	aDecimal, _ := strconv.ParseInt(aHex, 16, 0)
	bDecimal, _ := strconv.ParseInt(bHex, 16, 0)
	resultDecimal := aDecimal + bDecimal
	resultHex := strconv.FormatInt(resultDecimal, 16)
	return resultHex
}

func main() {
	pid, _ := bindDefaultProcess("winmine_(한글).exe")
	baseAddress, _ := memoryReadInit(pid)
	defer memoryReadClose()

	fmt.Println("Base address is", baseAddress)

	xPositionAddressInt, zPositionAddressInt := GetAddresses()
	fmt.Println("X address is", strconv.FormatInt(xPositionAddressInt, 16), "Z address is", strconv.FormatInt(zPositionAddressInt, 16))

	x := readMemoryAt(xPositionAddressInt)
	z := readMemoryAt(zPositionAddressInt)
	fmt.Println("X="+fmt.Sprintf("%g", x), "Z="+fmt.Sprintf("%g", z))
	var startAddress uintptr = 0x0000000000000000
	var endAddress uintptr = 0x7FFFFFFFFFFFFFFF

	// 메모리 영역 정보를 저장할 구조체
	var mbi windows.MemoryBasicInformation

	// 메모리 주소 범위 순회
	for addr := startAddress; ; {
		// 메모리 영역 정보 조회
		err := windows.VirtualQueryEx(Handle, addr, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			break // 더 이상 조회할 메모리 영역이 없으면 종료
		}

		// 메모리 영역 정보 출력
		fmt.Printf("BaseAddress: 0x%X, RegionSize: %d\n", mbi.BaseAddress, mbi.RegionSize)

		addr = mbi.BaseAddress + mbi.RegionSize
		endAddress = addr
	}

	var targetAddress uintptr = 0x0100579C

	// 특정 값을 찾기 위한 변수 설정
	var searchValue uint32 = 10 // 찾고자 하는 값

	// 메모리 순회
	for addr := startAddress; addr < endAddress; addr += unsafe.Sizeof(searchValue) {
		break

		var data uint32
		var bytesRead uintptr

		// 메모리 읽기
		err := windows.ReadProcessMemory(Handle, targetAddress, (*byte)(unsafe.Pointer(&data)), unsafe.Sizeof(data), &bytesRead)
		if err != nil {
			continue // 읽기 실패 시 다음 주소로
		}

		// 특정 값 찾기
		if data == searchValue {
			fmt.Printf("Found value at: 0x%X\n", addr)
			// 필요한 추가 작업 수행
		}
	}

	data, _ := readMemory(targetAddress)
	fmt.Println(data)

	if !isAdmin() {
		fmt.Println("This program needs to be run as administrator.")
		return
	}

	writeMemory(targetAddress, 42)
}

func writeMemory(address uintptr, data uint32) (bytesRead uintptr) {
	err := windows.WriteProcessMemory(Handle, address, (*byte)(unsafe.Pointer(&data)), unsafe.Sizeof(data), &bytesRead)
	if err != nil {
		return
	}

	return
}

type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func processes() ([]WindowsProcess, error) {
	const TH32CS_SNAPPROCESS = 0x00000002

	handle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	err = windows.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func findProcessByName(processes []WindowsProcess, name string) *WindowsProcess {
	for _, p := range processes {
		if strings.ToLower(p.Exe) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func newWindowsProcess(e *windows.ProcessEntry32) WindowsProcess {
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func bindDefaultProcess(defaultName string) (uint32, bool) {
	procs, err := processes()
	if err != nil {
		return 0, false
	}

	explorer := findProcessByName(procs, defaultName)
	if explorer == nil {
		return 0, false
	}

	return uint32(explorer.ProcessID), true
}

func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.GetCurrentProcessToken()
	defer token.Close()

	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}
