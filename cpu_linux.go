package sysinfo

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

func GetCPUInfo() any {
	totalPercent, _ := cpu.Percent(0, false)

	data := []struct {
		Identifier  string `json:"identifier"`
		Power       int    `json:"power"`
		Load        int    `json:"load"`
		Temperature int    `json:"temperature"`
	}{{"CPU", 0, int(totalPercent[0]), 0}}

	return data
}
