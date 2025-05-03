package dumper

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ReadProcessMemory(pid uint32, ea uintptr, size int) []byte {
	for region := range Regions(pid, windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ) {
		if ea < region.Metadata.BaseAddress || ea >= region.Metadata.BaseAddress+uintptr(region.Metadata.RegionSize) {
			continue
		}

		data, err := region.Read(ea, size)
		if err == nil {
			return data
		}
	}
	return nil
}

func ReadUInt32(pid uint32, ea uintptr) (uint32, error) {
	buffer := ReadProcessMemory(pid, ea, 4)
	if buffer == nil || len(buffer) < 4 {
		return 0, fmt.Errorf("failed to read 4 bytes at %x", ea)
	}
	return *(*uint32)(unsafe.Pointer(&buffer[0])), nil
}
