package sysinfo

func GetGPUInfo() any {
	data := []struct {
		Identifier  string `json:"identifier"`
		Power       int    `json:"power"`
		Load        int    `json:"load"`
		Temperature int    `json:"temperature"`
	}{}

	return data
}
