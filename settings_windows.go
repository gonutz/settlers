package main

import (
	"syscall"
	"unsafe"
)

// On Windows we need to call
//     SHGetFolderPath(0, CSIDL_APPDATA, NULL, SHGFP_TYPE_CURRENT, path)
// where path is of type LPTSTR. It is a ponter to a null-terminated string of
// length MAX_PATH which will receive the result. (declare TCHAR path[MAX_PATH])
// If the function succeeds, it will return S_OK.

func getSettingsPath() (path string) {
	path = "."
	shell32, err := syscall.LoadLibrary("Shell32.dll")
	if err != nil {
		return
	}
	getFolderPath, err := syscall.GetProcAddress(shell32, "SHGetFolderPathW")
	if err != nil {
		return
	}

	const (
		CSIDL_APPDATA      = 0x001A
		SHGFP_TYPE_CURRENT = 0
		S_OK               = 0
		MAX_PATH           = 260
	)

	var buffer [MAX_PATH]uint16
	pathPointer := uintptr(unsafe.Pointer(&buffer[0]))
	ret, _, err := syscall.Syscall6(
		uintptr(getFolderPath),
		5, // actual argument count
		0,
		CSIDL_APPDATA,
		0,
		SHGFP_TYPE_CURRENT,
		pathPointer,
		0,
	)
	if err != nil || ret != S_OK {
		println(err.Error())
		return
	}

	return syscall.UTF16ToString(buffer[:])
}
