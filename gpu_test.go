package sysinfo

import (
	"fmt"
	"sysinfo/wmi"
	"testing"
)

func TestWMI(t *testing.T) {
	var hms []openHardwareMonitorSensor
	err := wmi.Query("select * from Sensor", &hms, "127.0.0.1", "root/OpenHardwareMonitor")
	if err == nil {
		for _, item := range hms {
			fmt.Println(item)
		}
	}

	t.Errorf("%v ", "ans")
}

func TestGPU(t *testing.T) {
	ans := GetGPUInfo()

	t.Errorf("%v ", ans)
}
