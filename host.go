package sysinfo

import (
	"time"

	"github.com/shirou/gopsutil/v4/host"
	gopsutilnet "github.com/shirou/gopsutil/v4/net"
)

func GetBootTime() string {
	ts, _ := host.BootTime()
	data := time.Unix(int64(ts), 0).Format(time.DateTime)

	return data
}

func GetNetIO() any {
	type tmpStructA struct {
		BytesSent uint64 `json:"bytesSent"`
		BytesRecv uint64 `json:"bytesRecv"`
	}

	type tmpStructB struct {
		Ts   int64                 `json:"ts"`
		Data map[string]tmpStructA `json:"data"`
	}

	data := tmpStructB{
		Ts:   time.Now().Unix(),
		Data: map[string]tmpStructA{},
	}

	ios, _ := gopsutilnet.IOCounters(true)
	for _, item := range ios {
		data.Data[item.Name] = tmpStructA{BytesSent: item.BytesSent, BytesRecv: item.BytesRecv}
	}

	return data
}
