package win32

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	// BOOL ExitWindowsEx(UINT, DWORD)
	/*
		BOOL ExitWindowsEx(
		  [in] UINT  uFlags,
		  [in] DWORD dwReason
		);
	*/
	pExitWindowsEx = user32.NewProc("ExitWindowsEx")

	/*
		BOOL OpenProcessToken(
		  [in]  HANDLE  ProcessHandle,
		  [in]  DWORD   DesiredAccess,
		  [out] PHANDLE TokenHandle
		);
	*/
	pOpenProcessToken = advapi32.NewProc("OpenProcessToken")
	/*
		BOOL LookupPrivilegeValueW(
		  [in, optional] LPCWSTR lpSystemName,
		  [in]           LPCWSTR lpName,
		  [out]          PLUID   lpLuid
		);
	*/
	pLookupPrivilegeValue = advapi32.NewProc("LookupPrivilegeValueW")
	/*
		BOOL AdjustTokenPrivileges(
		  [in]            HANDLE            TokenHandle,
		  [in]            BOOL              DisableAllPrivileges,
		  [in, optional]  PTOKEN_PRIVILEGES NewState,
		  [in]            DWORD             BufferLength,
		  [out, optional] PTOKEN_PRIVILEGES PreviousState,
		  [out, optional] PDWORD            ReturnLength
		);
	*/
	pAdjustTokenPrivileges = advapi32.NewProc("AdjustTokenPrivileges")

	/*
		HANDLE GetCurrentProcess();
	*/
	pGetCurrentProcess = kernel32.NewProc("GetCurrentProcess")
	/*
		DWORD GetLastError();
	*/
	pGetLastError = kernel32.NewProc("GetLastError")
)

const (
	TOKEN_ADJUST_PRIVILEGES = 0x0020
	TOKEN_QUERY             = 0x0008

	SE_SHUTDOWN_NAME     = "SeShutdownPrivilege"
	SE_PRIVILEGE_ENABLED = 0x00000002
)

func GetCurrentProcess() uintptr {
	h, _, _ := pGetCurrentProcess.Call()

	return h
}

func OpenProcessToken(procH uintptr, access uint32, tokenH *uintptr) bool {
	rv, _, _ := pOpenProcessToken.Call(
		procH,                           // HANDLE
		uintptr(access),                 // DWORD
		uintptr(unsafe.Pointer(tokenH)), // PHANDLE
	)

	return rv > 0
}

func LookupPrivilegeValue(sysName string, name string, luid *Luid) bool {
	var (
		sysW  unsafe.Pointer
		nameW unsafe.Pointer
	)

	if len(sysName) > 0 {
		w, _ := windows.UTF16PtrFromString(sysName)
		sysW = unsafe.Pointer(w)
	}

	nameptr, _ := windows.UTF16PtrFromString(name)
	nameW = unsafe.Pointer(nameptr)

	rv, _, _ := pLookupPrivilegeValue.Call(
		uintptr(sysW),
		uintptr(nameW),
		uintptr(unsafe.Pointer(luid)),
	)

	return rv > 0
}

func AdjustTokenPrivileges(
	tokenH uintptr,
	disableAll bool,
	newState *TokenPrivileges,
	bufLen uint32,
	prevState *TokenPrivileges,
	rlen *uint32,
) bool {
	rv, _, _ := pAdjustTokenPrivileges.Call(
		tokenH,
		b2u(disableAll),
		uintptr(unsafe.Pointer(newState)),
		uintptr(bufLen),
		uintptr(unsafe.Pointer(prevState)),
		uintptr(unsafe.Pointer(rlen)),
	)

	return rv > 0
}

func GetLastError() error {
	return windows.GetLastError()
}

func b2u(b bool) uintptr {
	if b {
		return 1
	}

	return 0
}

const (
	EWX_LOGOFF      = 0
	EWX_SHUTDOWN    = 1
	EWX_REBOOT      = 2
	EWX_FORCE       = 4
	EWX_POWEROFF    = 8
	EWX_FORCEIFHUNG = 16

	SHTDN_REASON_FLAG_PLANNED          = 0x80000000
	SHTDN_REASON_MAJOR_OPERATINGSYSTEM = 0x00020000
	SHTDN_REASON_MINOR_UPGRADE         = 0x00000003
)

func ExitWindowsEx(uflags uint, reason uint32) bool {
	rv, _, _ := pExitWindowsEx.Call(
		uintptr(uflags),
		uintptr(reason),
	)

	return rv > 0
}

/*
typedef struct _TOKEN_PRIVILEGES {
  DWORD               PrivilegeCount;
  LUID_AND_ATTRIBUTES Privileges[ANYSIZE_ARRAY]; // 1 in pub header
} TOKEN_PRIVILEGES, *PTOKEN_PRIVILEGES;
*/

type TokenPrivileges struct {
	PrivilegeCount uint32
	Privileges     [1]LuidAndAttributes
}

/*
	typedef struct _LUID {
	  DWORD LowPart;
	  LONG  HighPart;
	} LUID, *PLUID;
*/
type Luid struct {
	LowPart  uint32 // DWORD
	HighPart int32  // LONG, maybe actually 64?
}

/*
	typedef struct _LUID_AND_ATTRIBUTES {
	  LUID  Luid;
	  DWORD Attributes;
	} LUID_AND_ATTRIBUTES, *PLUID_AND_ATTRIBUTES;
*/
type LuidAndAttributes struct {
	Luid       Luid
	Attributes uint32 // DWORD
}
