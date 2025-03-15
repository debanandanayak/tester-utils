package executable

import "syscall"

func killProcess(pid int) {
	syscall.Kill(pid, syscall.SIGTERM)  // Don't know if this is required
	syscall.Kill(-pid, syscall.SIGTERM) // Kill the whole process group
}

func createProcAttribute() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}
