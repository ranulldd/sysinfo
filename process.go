package sysinfo

import (
	"slices"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

func GetProcessInfo(pList []string) any {
	type tmpStruct struct {
		Pid        int32  `json:"pid"`
		Name       string `json:"name"`
		CPU        int    `json:"cpu"`
		CreateTime string `json:"createTime"`
		RSS        uint64 `json:"rss"`
		VMS        uint64 `json:"vms"`
		ReadBytes  uint64 `json:"readBytes"`
		WriteBytes uint64 `json:"writeBytes"`
	}

	data := []tmpStruct{}

	ps, _ := process.Processes()
	for _, p := range ps {
		name, _ := p.Name()
		if slices.Contains(pList, name) {
			CPUPercent, _ := p.CPUPercent()
			CreateTime, _ := p.CreateTime()
			memInfo, _ := p.MemoryInfo()
			rss := memInfo.RSS
			vms := memInfo.VMS
			ioInfo, _ := p.IOCounters()
			ReadBytes := ioInfo.ReadBytes
			WriteBytes := ioInfo.WriteBytes

			data = append(data, tmpStruct{
				Pid:        p.Pid,
				Name:       name,
				CPU:        int(CPUPercent),
				CreateTime: time.Unix(CreateTime/1000, 0).Format(time.DateTime),
				RSS:        rss,
				VMS:        vms,
				ReadBytes:  ReadBytes,
				WriteBytes: WriteBytes,
			})
		}

	}

	return data
}
