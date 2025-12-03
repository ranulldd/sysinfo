package sysinfo

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/ranulldd/sysinfo/wmi"
)

// OpenHardwareMonitor needed

func GetGPUInfo() any {
	type dataStruct struct {
		Identifier  string `json:"identifier"`
		Power       int    `json:"power"`
		Load        int    `json:"load"`
		Temperature int    `json:"temperature"`
	}
	data := []dataStruct{}
	var hms []openHardwareMonitorSensor
	err := wmi.Query("select * from Sensor where name='GPU Core' or name='GPU Total'", &hms, "127.0.0.1", "root/OpenHardwareMonitor")
	if err == nil {
		for _, item := range hms {
			idx := 0
			found := false
			for i, v := range data {
				if v.Identifier == item.Parent {
					idx = i
					found = true
					break
				}
			}
			if !found {
				data = append(data, dataStruct{Identifier: item.Parent})
				idx = len(data) - 1
			}

			if item.Name == "GPU Total" && item.SensorType == "Power" {
				data[idx].Power = int(item.Value)
			} else if item.Name == "GPU Core" && item.SensorType == "Load" {
				data[idx].Load = int(item.Value)
			} else if item.Name == "GPU Core" && item.SensorType == "Temperature" {
				data[idx].Temperature = int(item.Value)
			}
		}
	}

	cmd := exec.Command("nvidia-smi", "--query-gpu=name,power.draw,utilization.gpu,temperature.gpu", "--format=csv,noheader,nounits")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return data
	}
	items := strings.Split(string(output), ",")
	if len(items) != 4 || !strings.Contains(string(output), "Tesla") {
		return data
	}

	data = append(data, dataStruct{Identifier: items[0]})
	idx := len(data) - 1

	v, _ := strconv.ParseFloat(strings.TrimSpace(items[1]), 64)
	data[idx].Power = int(v)
	v, _ = strconv.ParseFloat(strings.TrimSpace(items[2]), 64)
	data[idx].Load = int(v)
	v, _ = strconv.ParseFloat(strings.TrimSpace(items[3]), 64)
	data[idx].Temperature = int(v)

	return data
}
