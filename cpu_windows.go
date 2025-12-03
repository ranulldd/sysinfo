package sysinfo

import "github.com/ranulldd/sysinfo/wmi"

// OpenHardwareMonitor needed

// type openHardwareMonitorHardware struct {
// 	Name         string
// 	Identifier   string
// 	HardwareType string
// 	Parent       string
// }

type openHardwareMonitorSensor struct {
	Name       string
	Value      float32
	SensorType string
	Identifier string
	Parent     string
}

func GetCPUInfo() any {
	type dataStruct struct {
		Identifier  string `json:"identifier"`
		Power       int    `json:"power"`
		Load        int    `json:"load"`
		Temperature int    `json:"temperature"`
	}
	data := []dataStruct{}
	var hms []openHardwareMonitorSensor
	err := wmi.Query("select * from Sensor where name='CPU Package' or name='CPU Total'", &hms, "127.0.0.1", "root/OpenHardwareMonitor")
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

			if item.Name == "CPU Package" && item.SensorType == "Power" {
				data[idx].Power = int(item.Value)
			} else if item.Name == "CPU Total" && item.SensorType == "Load" {
				data[idx].Load = int(item.Value)
			} else if item.Name == "CPU Package" && item.SensorType == "Temperature" {
				data[idx].Temperature = int(item.Value)
			}
		}
	}

	return data
}
