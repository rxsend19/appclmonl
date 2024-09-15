package foreground

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/abdfnx/gosh"
	"golang.org/x/sys/windows"
)

// Иногда сыпет названием в bytes с недекодируемым ASCII
func getWindowThreadProcessId(hwnd uintptr) uintptr {
	var prcsId uintptr = 0
	us32 := syscall.MustLoadDLL("user32.dll")
	prc := us32.MustFindProc("GetWindowThreadProcessId")
	prc.Call(hwnd, uintptr(unsafe.Pointer(&prcsId)))
	return prcsId
}

func Run() string {
	HWND := windows.GetForegroundWindow()
	_, Raw, _ := gosh.RunOutput(`tasklist /FI "pid eq ` + fmt.Sprint(getWindowThreadProcessId(uintptr(HWND))) + `" /FO table /NH`)
	exe := strings.TrimSpace(strings.Split(Raw, " ")[0])

	return exe
}
