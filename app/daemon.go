package app

import (
	"flag"
	"os"
	"path/filepath"
	"scaffold/common/config"

	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

// service manager
var serviceType = flag.String("s", "", "Services Management, install, uninstall")

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go start()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func getService() service.Service {
	options := make(service.KeyValue)
	svcConfig := &service.Config{
		Name:        config.Config.Service.Name,
		DisplayName: config.Config.Service.DisplayName,
		Description: config.Config.Service.Description,
		Option:      options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logrus.Fatalln(err)
	}
	return s
}

func installService() {
	s := getService()
	status, err := s.Status()
	if err != nil && status == service.StatusUnknown {
		// 服务未知，创建服务
		if err = s.Install(); err == nil {
			s.Start()
			logrus.Printf("install %s install successful!", config.Config.Service.Name)
			return
		}
		logrus.Printf("install %s service failure, ERR: %s\n", err, config.Config.Service.Name)
	}

	if status != service.StatusUnknown {
		logrus.Printf("%s service installed, no reinstallation required", config.Config.Service.Name)
	}
}

func uninstallService() {
	s := getService()
	s.Stop()
	if err := s.Uninstall(); err == nil {
		logrus.Printf("%s service uninstall successful!", config.Config.Service.Name)
	} else {
		logrus.Printf("%s service uninstall failure, ERR: %s\n", err, config.Config.Service.Name)
	}
	os.Exit(1)
}

func StartDaemon() {
	flag.Parse()
	switch *serviceType {
	case "install":
		installService()
	case "uninstall":
		uninstallService()
	default:
		s := getService()
		status, _ := s.Status()
		if status != service.StatusUnknown {
			setCurrentDirToExecutableDir()
			// service runs
			s.Run()
		} else {
			logrus.Println("non-service runs")
			switch s.Platform() {
			case "windows-service":
				logrus.Printf("service runs: .\\%s.exe -s install", config.Config.Service.Name)
			default:
				logrus.Printf("service runs: sudo ./%s -s install", config.Config.Service.Name)
			}
			// run anything
			s.Run()
		}
	}
}

// setCurrentDirToExecutableDir 设置当前工作目录为可执行文件所在的目录
func setCurrentDirToExecutableDir() error {
	// 获取可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// 获取可执行文件的目录
	exeDir := filepath.Dir(exePath)

	// 将当前工作目录设置为可执行文件的目录
	err = os.Chdir(exeDir)
	if err != nil {
		return err
	}

	return nil
}
