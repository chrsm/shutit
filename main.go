package main

import (
	"log"
	"runtime"
	"time"

	"bits.chrsm.org/shutit/internal/win32"
)

func main() {
	runtime.LockOSThread()

	tick := time.NewTicker(30 * time.Second)
	for {
		select {
		case t := <-tick.C:
			if t.Local().Hour() > 22 || t.Local().Hour() < 7 {
				log.Println("Get off the fuckin' computer")
				shutit()

				return
			}
		}
	}
}

func shutit() {
	var (
		token uintptr
		privs win32.TokenPrivileges
	)

	hnd := win32.GetCurrentProcess()
	if !win32.OpenProcessToken(hnd, win32.TOKEN_ADJUST_PRIVILEGES|win32.TOKEN_QUERY, &token) {
		log.Printf("couldn't open current process; token=%#v, hnd=%#v", token, hnd)
		return
	}

	win32.LookupPrivilegeValue(
		"",
		win32.SE_SHUTDOWN_NAME,
		&privs.Privileges[0].Luid,
	)

	privs.PrivilegeCount = 1
	privs.Privileges[0].Attributes = win32.SE_PRIVILEGE_ENABLED

	if !win32.AdjustTokenPrivileges(token, false, &privs, 0, nil, nil) {
		log.Println("failed to adjust token privs")
	}

	err := win32.GetLastError()
	if err != nil {
		log.Printf("GetLastError: %#v", err)
		return
	}

	ok := win32.ExitWindowsEx(
		win32.EWX_SHUTDOWN|win32.EWX_FORCE,
		win32.SHTDN_REASON_FLAG_PLANNED|win32.SHTDN_REASON_MINOR_UPGRADE|win32.SHTDN_REASON_MAJOR_OPERATINGSYSTEM,
	)
	if !ok {
		log.Printf("failed to shut down: %t", ok)
	}
}
