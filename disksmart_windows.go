package sysinfo

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/ranulldd/sysinfo/wmi"
)

type diskInfoStruct struct {
	DeviceID    string
	Model       string
	Temperature int
	Life        string
	Remark      string
}

type win32_DiskDriveStruct struct {
	DeviceID     string
	Size         string
	Model        string
	SerialNumber string
	PnpDeviceId  string
}

type msStorageDriver_FailurePredictDataStruct struct {
	InstanceName   string
	VendorSpecific []byte
}

type storageProtocolSpecificDataStruct struct { // 40 bytes
	ProtocolType                uint32
	DataType                    uint32
	ProtocolDataRequestValue    uint32
	ProtocolDataRequestSubValue uint32
	ProtocolDataOffset          uint32
	ProtocolDataLength          uint32
	FixedProtocolReturnData     uint32
	Reserved                    [3]uint32
}

type storagePropertyQueryStruct struct { // 8 bytes
	PropertyId uint32
	QueryType  uint32
}

type storageQueryWithBufferStruct struct {
	Query            storagePropertyQueryStruct        // 8 bytes
	ProtocolSpecific storageProtocolSpecificDataStruct // 40 bytes
	Buffer           [4096]byte
}

func fillNVMESmartData(diskInfo *diskInfoStruct) {
	name, _ := syscall.UTF16FromString(diskInfo.DeviceID)
	fd, err := syscall.CreateFile(
		&name[0],
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0)
	if err != nil {
		return
	}
	defer syscall.CloseHandle(fd)

	var nptwb storageQueryWithBufferStruct
	var returnLen uint32
	nptwb.ProtocolSpecific.ProtocolType = 3
	nptwb.ProtocolSpecific.DataType = 2
	nptwb.ProtocolSpecific.ProtocolDataRequestValue = 2
	nptwb.ProtocolSpecific.ProtocolDataRequestSubValue = 0
	nptwb.ProtocolSpecific.ProtocolDataOffset = uint32(unsafe.Sizeof(nptwb.ProtocolSpecific))
	nptwb.ProtocolSpecific.ProtocolDataLength = 4096
	nptwb.Query.PropertyId = 49
	nptwb.Query.QueryType = 0

	err = syscall.DeviceIoControl(
		fd,
		((0x0000002d)<<16)|((0)<<14)|((0x0500)<<2), // IOCTL_STORAGE_QUERY_PROPERTY
		(*byte)(unsafe.Pointer(&nptwb)),
		uint32(unsafe.Sizeof(nptwb)),
		(*byte)(unsafe.Pointer(&nptwb)),
		uint32(unsafe.Sizeof(nptwb)),
		&returnLen,
		nil)
	if err != nil {
		nptwb.ProtocolSpecific.ProtocolDataRequestSubValue = 0xFFFFFFFF
		err = syscall.DeviceIoControl(
			fd,
			((0x0000002d)<<16)|((0)<<14)|((0x0500)<<2), // IOCTL_STORAGE_QUERY_PROPERTY
			(*byte)(unsafe.Pointer(&nptwb)),
			uint32(unsafe.Sizeof(nptwb)),
			(*byte)(unsafe.Pointer(&nptwb)),
			uint32(unsafe.Sizeof(nptwb)),
			&returnLen,
			nil)
	}
	if err != nil {
		return
	}

	diskInfo.Temperature = int(nptwb.Buffer[0x2])*256 + int(nptwb.Buffer[0x1]) - 273
	diskInfo.Life = fmt.Sprintf("%v%%", 100-nptwb.Buffer[0x05])
	// log.Printf("%v", nptwb.Buffer)
}

func fillATASmartData(diskInfos []*diskInfoStruct, diskDrives []win32_DiskDriveStruct, msStorageDriver msStorageDriver_FailurePredictDataStruct) {
	msStorageDriver.InstanceName = strings.ToUpper(msStorageDriver.InstanceName)
	for diskId, diskDrive := range diskDrives {
		if strings.Contains(msStorageDriver.InstanceName, diskDrive.PnpDeviceId) {
			for i := 0; i < 30; i++ {
				idx := 2 + i*12
				switch msStorageDriver.VendorSpecific[idx] {
				case 0x01: // Read Error Rate
					readErrorRate := msStorageDriver.VendorSpecific[idx+3]
					if readErrorRate > 100 {
						if readErrorRate < 200 {
							diskInfos[diskId].Remark = "出现读取错误"
							diskInfos[diskId].Life = "警告"
						}
					}
				case 0xC2: // Temperature
					diskInfos[diskId].Temperature = int(msStorageDriver.VendorSpecific[idx+5])
				case 0xCA: // Lifetime Remaining
					diskInfos[diskId].Life = fmt.Sprintf("%v%%", msStorageDriver.VendorSpecific[idx+3])
				}
			}

			break
		}
	}
}

func getDiskSmart() []*diskInfoStruct {
	var diskDrives []win32_DiskDriveStruct
	err := wmi.Query("select * from Win32_DiskDrive", &diskDrives)
	if err != nil {
		return nil
	}

	diskInfos := make([]*diskInfoStruct, len(diskDrives))
	for i, diskDrive := range diskDrives {
		diskInfos[i] = &diskInfoStruct{
			DeviceID: diskDrive.DeviceID,
			Model:    diskDrive.Model,
			Life:     "良好",
		}
		if strings.Contains(diskDrive.Model, "NVMe") || strings.Contains(diskDrive.Model, "Optane") || strings.Contains(diskDrive.PnpDeviceId, "NVME") || strings.Contains(diskDrive.PnpDeviceId, "OPTANE") {
			fillNVMESmartData(diskInfos[i])
		}
	}

	var msStorageDriver_FailurePredictDatas []msStorageDriver_FailurePredictDataStruct

	err = wmi.Query("select * from MSStorageDriver_FailurePredictData", &msStorageDriver_FailurePredictDatas, "127.0.0.1", "root/wmi")
	if err != nil {
		return nil
	}

	for _, msStorageDriver := range msStorageDriver_FailurePredictDatas {
		fillATASmartData(diskInfos, diskDrives, msStorageDriver)
	}

	return diskInfos
}

func GetDiskSmartInfo() any {
	type tmpStruct struct {
		Model       string `json:"model"`
		Temperature int    `json:"temperature"`
		Life        string `json:"life"`
	}

	data := []tmpStruct{}

	diskInfos := getDiskSmart()
	if len(diskInfos) > 0 {
		data = make([]tmpStruct, len(diskInfos))
		for i, diskInfo := range diskInfos {
			data[i] = tmpStruct{diskInfo.Model, diskInfo.Temperature, diskInfo.Life}
		}
	}

	return data
}
