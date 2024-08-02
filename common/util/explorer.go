package util

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
)

// DockerEnvFile Docker容器中包含的文件
const DockerEnvFile string = "/.dockerenv"

// IsRunInDocker 是否在docker中运行
func IsRunInDocker() bool {
	_, err := os.Stat(DockerEnvFile)
	return err == nil
}

func AutoOpenExplorer(listen string) {
	if IsRunInDocker() {
		// docker中运行, 提示
		fmt.Println("runing in docker")
	} else {
		addr, err := net.ResolveTCPAddr("tcp", listen)
		if err != nil {
			return
		}
		url := fmt.Sprintf("http://127.0.0.1:%d", addr.Port)
		if addr.IP.IsGlobalUnicast() {
			url = fmt.Sprintf("http://%s", addr.String())
		}
		go openExplorer(url)
	}
}

// OpenExplorer Open local browser
func openExplorer(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		// mac
		cmd = "open"
	default:
		// linux
		cmd = "xdg-open"
	}
	args = append(args, url)

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Printf("Please open the browser manually and visit %s", url)
	} else {
		fmt.Println("Browser opened successfully")
	}
}
