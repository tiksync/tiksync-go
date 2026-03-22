package tiksync

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

//go:embed native/tiksync_core.dll
var embeddedDLL embed.FS

type dynamicNativeCore struct {
	dll        *syscall.DLL
	initProc   *syscall.Proc
	freeProc   *syscall.Proc
	headersProc *syscall.Proc
	freeStrProc *syscall.Proc
	ptr        uintptr
}

func newDynamicNativeCore(apiKey string) *dynamicNativeCore {
	if runtime.GOOS != "windows" {
		return nil
	}

	dllPath, err := extractDLL()
	if err != nil {
		return nil
	}

	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		return nil
	}

	initProc, _ := dll.FindProc("tiksync_init")
	freeProc, _ := dll.FindProc("tiksync_free")
	headersProc, _ := dll.FindProc("tiksync_get_headers_json")
	freeStrProc, _ := dll.FindProc("tiksync_free_string")

	if initProc == nil || freeProc == nil || headersProc == nil || freeStrProc == nil {
		dll.Release()
		return nil
	}

	keyBytes := []byte(apiKey)
	ptr, _, _ := initProc.Call(
		uintptr(unsafe.Pointer(&keyBytes[0])),
		uintptr(len(keyBytes)),
	)
	if ptr == 0 {
		dll.Release()
		return nil
	}

	return &dynamicNativeCore{
		dll:         dll,
		initProc:    initProc,
		freeProc:    freeProc,
		headersProc: headersProc,
		freeStrProc: freeStrProc,
		ptr:         ptr,
	}
}

func (d *dynamicNativeCore) getHeaders() map[string]string {
	if d == nil || d.ptr == 0 {
		return nil
	}

	jsonPtr, _, _ := d.headersProc.Call(d.ptr)
	if jsonPtr == 0 {
		return nil
	}

	jsonStr := ptrToString(jsonPtr)
	d.freeStrProc.Call(jsonPtr)

	var headers map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &headers); err != nil {
		return nil
	}
	return headers
}

func (d *dynamicNativeCore) close() {
	if d != nil && d.ptr != 0 {
		d.freeProc.Call(d.ptr)
		d.ptr = 0
	}
	if d != nil && d.dll != nil {
		d.dll.Release()
	}
}

func extractDLL() (string, error) {
	tmpDir := filepath.Join(os.TempDir(), "tiksync-native")
	os.MkdirAll(tmpDir, 0755)
	dllPath := filepath.Join(tmpDir, "tiksync_core.dll")

	if _, err := os.Stat(dllPath); err == nil {
		return dllPath, nil
	}

	data, err := embeddedDLL.ReadFile("native/tiksync_core.dll")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(dllPath, data, 0755); err != nil {
		return "", err
	}

	return dllPath, nil
}

func ptrToString(ptr uintptr) string {
	var bytes []byte
	for {
		b := *(*byte)(unsafe.Pointer(ptr))
		if b == 0 {
			break
		}
		bytes = append(bytes, b)
		ptr++
	}
	return string(bytes)
}
