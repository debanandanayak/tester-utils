package executable

import (
	"log"
	"syscall"
)

func killProcess(pid int) {
	handle, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		log.Fatalf("Failed to open process: %v", err)
	}
	defer syscall.CloseHandle(handle)

	err = syscall.TerminateProcess(handle, 1) // Exit code 1
	if err != nil {
		log.Fatalf("Failed to terminate process: %v", err)
	}

	log.Println("Process terminated successfully.")
}

func createProcAttribute() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}
