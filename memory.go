package sysinfo

import "github.com/shirou/gopsutil/v4/mem"

func GetMemoryInfo() any {

	var data struct {
		Used  uint64 `json:"used"`
		Total uint64 `json:"total"`
	}

	v, _ := mem.VirtualMemory()
	data.Used, data.Total = v.Used, v.Total

	return data
}
