package process

import (
	"errors"
	"strings"
	"syscall"
	"unsafe"

	"github.com/0xrawsec/golang-win32/win32/kernel32"
	"golang.org/x/sys/windows"
)

var Handle windows.Handle

func Open(pid uint32) (windows.Handle, error) {
	handle, err := windows.OpenProcess(kernel32.PROCESS_ALL_ACCESS, false, pid)
	if err != nil {
		return 0, err
	}

	return handle, nil
}

func Close() {
	windows.CloseHandle(Handle)
}

type process struct {
	ProcessID       int
	ParentProcessID int
	Name            string
}

func getProcessList() ([]process, error) {
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

	results := make([]process, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			if errors.Is(err, syscall.ERROR_NO_MORE_FILES) {
				return results, nil
			}
			return nil, err
		}
	}
}

func findProcessByName(processes []process, name string) *process {
	for _, p := range processes {
		if strings.ToLower(p.Name) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func newWindowsProcess(e *windows.ProcessEntry32) process {
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return process{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Name:            syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func GetPid(processName string) (pid uint32, found bool) {
	processList, err := getProcessList()
	if err != nil {
		return 0, false
	}

	explorer := findProcessByName(processList, processName)
	if explorer == nil {
		return 0, false
	}

	return uint32(explorer.ProcessID), true
}

func HasAdminPermission() bool {
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
