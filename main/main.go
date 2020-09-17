package main

import (
	"cbinstaller/core"
	"cbinstaller/util"
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"os"
)

func main() {

	input := os.Args[1:]

	var command string
	if len(input) == 0 {
		fmt.Println("\nCBwindowInstaller.exe <command> [-parameter]")
		fmt.Println("\nYou can verify command and description by using command 'CBwindowInstaller.exe --help'")
		return
	} else {
		command = input[0]
		if command != "--help" && command != "install" && command != "uninstall" && command != "start" && command != "stop" && command != "restart" {
			fmt.Println("\nCBwindowInstaller.exe <command> [-parameter]")
			fmt.Println("\nYou can verify command and description by using command 'CBwindowInstaller.exe --help'")
			return
		}
	}
	os.Args = os.Args[1:]
	installPath := "C:/Program Files/telegraf"
	serviceName := "CB-Dragonfly agent"
	s := core.GetServiceInstance(
		&service.Config{
			Name: serviceName,
		})

	resourceInfo := core.GetResourceInstance(&core.ResourceInfo{
		TelegrafGetUrl: "/telegraf",
		ConfGetUrl:     "/conf",
		ServerInfo:     "192.168.130.7",
		InstallPath:    installPath,
	})
	flag.Usage = func() {
		fmt.Println("cbinstaller.exe provide beblow command lines.\n")
		fmt.Println("	Install monitoring agent (vmID, mcisID, cspType must required)\n")
		fmt.Println("		cbinstaller.exe install -namespace {{namepace}} -mcisID {{MCISID}} -cspType {{CSPID}} -vmID {{HOSTID}}\n")
		fmt.Println("	Uninstall monitoring agent\n")
		fmt.Println("		cbinstaller.exe uninstall\n")
		fmt.Println("	Start monitoring agent\n")
		fmt.Println("		cbinstaller.exe start\n")
		fmt.Println("	Stop monitoring agent\n")
		fmt.Println("		cbinstaller.exe stop\n")
		fmt.Println("	Restart monitoring agent\n")
		fmt.Println("		cbinstaller.exe restart\n")
		fmt.Println("Monitoring agent conf file directory : C:/Program Files/telegraf/telegraf.conf")
	}
	checkInstallParam := resourceInfo.CBInstallFlag(command, resourceInfo)

	switch command {
	case "install":
		if !checkInstallParam {
			return
		}
		CBAgentExeExist := util.CheckFileExists(resourceInfo.InstallPath + "/telegraf.exe")
		CBAgentConfExist := util.CheckFileExists("C:/Program Files/telegraf/telegraf.conf")
		if !CBAgentExeExist || !CBAgentConfExist {
			err := resourceInfo.InstallTelegrafwithConf(!CBAgentExeExist, !CBAgentConfExist)
			if err != nil {
				fmt.Println("\nFail to Download telegraf from server")
				resourceInfo.UnInstallTelegrafwithConf()
				return
			}
		}
		s = core.GetServiceInstance(
			&service.Config{
				Name:        serviceName,
				DisplayName: serviceName,
				Description: "Monitoring host resource",
				Executable:  installPath + "/telegraf.exe",
			})
		err := s.Install()
		if !checkInstallParam {
			return
		}
		if err != nil {
			fmt.Println("\nCB-Dragonfly Agent is already installed")
			return
		}
		_ = s.Start()
		fmt.Println("\nCB-Dragonfly Agent in installed")
		break

	case "uninstall":
		_ = s.Stop()
		err := s.Uninstall()
		// TODO: REMOVE EXE FILE AND CONF
		if err != nil {
			fmt.Println("\nCB-Dragonfly Agent in not installed")
			resourceInfo.UnInstallTelegrafwithConf()
			return
		}
		resourceInfo.UnInstallTelegrafwithConf()
		fmt.Println("\nCB-Dragonfly Agent is un-installed")
		break

	case "start":
		err := s.Start()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nCB-Dragonfly Agent Service Start")
		break

	case "stop":
		err := s.Stop()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nCB-Dragonfly Agent Service Stop")
		break

	case "restart":
		err := s.Restart()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nCB-Dragonfly Agent Restart")
		break
	}
}
