package api

import (
	"errors"

	"bits.chrsm.org/shutit/internal/win32"
)

func Shutdown() error {
	var (
		token uintptr
		privs win32.TokenPrivileges
	)

	hnd := win32.GetCurrentProcess()
	if !win32.OpenProcessToken(hnd, win32.TOKEN_ADJUST_PRIVILEGES|win32.TOKEN_QUERY, &token) {
		//log.Printf("couldn't open current process; token=%#v, hnd=%#v", token, hnd)
		return errors.New("couldn't open current process")
	}

	win32.LookupPrivilegeValue(
		"",
		win32.SE_SHUTDOWN_NAME,
		&privs.Privileges[0].Luid,
	)

	privs.PrivilegeCount = 1
	privs.Privileges[0].Attributes = win32.SE_PRIVILEGE_ENABLED

	if !win32.AdjustTokenPrivileges(token, false, &privs, 0, nil, nil) {
		//log.Println("failed to adjust token privs")
		return errors.New("failed to adjust token privs")
	}

	err := win32.GetLastError()
	if err != nil {
		//log.Printf("GetLastError: %#v", err)
		return errors.New("could not GetLastError")
	}

	ok := win32.ExitWindowsEx(
		win32.EWX_SHUTDOWN|win32.EWX_FORCE,
		win32.SHTDN_REASON_FLAG_PLANNED|win32.SHTDN_REASON_MINOR_UPGRADE|win32.SHTDN_REASON_MAJOR_OPERATINGSYSTEM,
	)
	if !ok {
		//log.Printf("failed to shut down: %t", ok)
		return errors.New("could not call ExitWindowsEx")
	}

	return nil
}
