package app

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"scaffold/internal/config"

	"github.com/kardianos/service"
)

// service manager
var serviceType = flag.String("s", "", "Services Management, install, uninstall")

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go run()
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
		slog.Error(err.Error())
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
			slog.Info("install service successful!")
			return
		}
		slog.Error("install service failure")
	}

	if status != service.StatusUnknown {
		slog.Info("service installed, no reinstallation required")
	}
}

func uninstallService() {
	s := getService()
	s.Stop()
	if err := s.Uninstall(); err == nil {
		slog.Info("service uninstall successful!")
	} else {
		slog.Error("service uninstall failure!")
	}
	os.Exit(1)
}

func startDaemon() {
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
			slog.Info("non-service runs")
			switch s.Platform() {
			case "windows-service":
				slog.Info(fmt.Sprintf("service runs: .\\%s.exe -s install", config.Config.Service.Name))
			default:
				slog.Info(fmt.Sprintf("service runs: sudo ./%s -s install", config.Config.Service.Name))
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
