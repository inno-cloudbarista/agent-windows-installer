package core

import (
	"cbinstaller/util"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type ResourceInfo struct {
	TelegrafGetUrl string
	ConfGetUrl string
	ServerInfo string
	InstallPath string
	ConfigMap map[string] string
}

func GetResourceInstance(rscInfo *ResourceInfo) *ResourceInfo {
	return rscInfo
}

func (r *ResourceInfo) CBInstallFlag(command string, rscInfo *ResourceInfo) bool {
	vmID := flag.String("vmID", "", "Input vmID")                // 명령줄 옵션을 받은 뒤 문자열로 저장
	mcisID := flag.String("mcisID", "", "Input mcisID")          // 명령줄 옵션을 받은 뒤 정수로 저장
	cspType := flag.String("cspType", "", "Input cspType")       // 명령줄 옵션을 받은 뒤 실수로 저장
	namespace := flag.String("namespace", "", "Input namespace") // 명령줄 옵션을 받은 뒤 실수로 저장
	flag.Parse()
	nFlag := flag.NFlag()

	if command == "--help" {
		flag.Usage()
		return false
	} else if command == "install" {
		if nFlag != 4 {
			flag.Usage()
			return false
		}

		configMap := map[string]string{}
		configMap["vm_id"] = *vmID
		configMap["mcis_id"] = *mcisID
		configMap["csp_type"] = *cspType
		configMap["ns_id"] = *namespace

		rscInfo.ConfigMap = configMap
		return true
	}
	return true
}

func (r *ResourceInfo) InstallTelegrafwithConf(installExe bool, installConf bool) error {
	if installConf {
		fmt.Println("\nTry to create telegraf.conf file\n")
		err := r.GetTelegrafConfFromServer()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("\nCompelete telegraf.exe creation. (C:/Program Files/telegraf/telegraf.conf)\n")
	}
	if installExe {
		fmt.Println("\nTry to install telegraf.exe\n")
		err := r.GetTelegrafExeFromServer()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("\nCompelete telegraf.exe installation\n")
		time.Sleep(1*time.Second)
	}
	return nil
}

func (r *ResourceInfo) GetTelegrafExeFromServer() error {
	var arch string
	is64Bit := uint64(^uintptr(0)) == ^uint64(0)
	if is64Bit {
		arch = "64"
	} else {
		arch = "32"
	}
	url := fmt.Sprintf("http://%s:9090/dragonfly/file/agent/pkg?osType=windows&arch=%s", r.ServerInfo, arch)
	installPath := r.InstallPath
	err := util.GetZipFileByApiResponse(installPath,url)
	if err != nil {
		return  err
	}
	err = util.GetUnZipFile(installPath)
	if err != nil {
		return  err
	}
	return nil
}

func (r *ResourceInfo) GetTelegrafConfFromServer() error {
	url := fmt.Sprintf("http://%s:9090/dragonfly/file/agent/conf?ns_id=%s&mcis_id=%s&vm_id=%s&csp_type=%s", r.ServerInfo, r.ConfigMap["ns_id"], r.ConfigMap["mcis_id"], r.ConfigMap["vm_id"], r.ConfigMap["csp_type"])
	byteConf, err := util.GetApiResponse(url)
	if err != nil {
		return  err
	}
	conf := string(byteConf)
	conf = strings.ReplaceAll(conf, `osType = "linux"`, `osType = "windows"`)
	conf = strings.ReplaceAll(conf, "{{influxdb_server}}", fmt.Sprintf("http://%s:8086", r.ServerInfo))
	// TODO : 카프카 적용시 아래 라인 변경 요망
	conf = strings.ReplaceAll(conf, "{{collector_server}}", fmt.Sprintf("udp://%s:8094", r.ServerInfo))

	_ = os.Mkdir("C:/Program Files/telegraf", 0755)
	f, err := os.Create("C:/Program Files/telegraf/telegraf.conf")
	if err != nil {
		return  err
	}

	_, err = f.WriteString(conf)
	if err != nil {
		return  err
	}
	return nil
}

func (r *ResourceInfo) UnInstallTelegrafwithConf() {
	_ = os.Remove(r.InstallPath+"/telegraf.exe")
	_ = os.RemoveAll("C:/Program Files/telegraf")
}