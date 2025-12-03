package sysinfo

func GetDiskSmartInfo() any {
	type tmpStruct struct {
		Model       string `json:"model"`
		Temperature int    `json:"temperature"`
		Life        string `json:"life"`
	}

	data := []tmpStruct{}

	return data
}
