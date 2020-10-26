package utils

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	DefaultDirMode os.FileMode = 0755
)

func ExecDir() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	path, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	splitstr := strings.Split(path, "/")
	return strings.Join(splitstr[:len(splitstr)-1], "/")
}

func IsDirExist(path string) bool {
	dir, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return dir.IsDir()
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
