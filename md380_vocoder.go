package md380vocoder

/*
#cgo CFLAGS: -I${SRCDIR}/go_wrapper/include
#cgo LDFLAGS: -L${SRCDIR}/go_wrapper/lib -lmd380_vocoder -lstdc++ -lm
#cgo linux LDFLAGS: -ldl -lpthread
#cgo darwin LDFLAGS: -framework CoreFoundation

#include "md380_vocoder_c.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

const (
	PCMFrameSize  = 160 // 160 samples per frame
	AMBEFrameSize = 9   // 9 bytes per AMBE frame
)

// Vocoder MD380 vocoder 实例
type Vocoder struct {
	handle C.md380_vocoder_handle_t
}

// NewVocoder 创建新的 vocoder 实例
func NewVocoder() (*Vocoder, error) {
	handle := C.md380_vocoder_create()
	if handle == nil {
		return nil, errors.New("failed to create vocoder: out of memory")
	}

	v := &Vocoder{handle: handle}

	// 设置终结器自动清理
	runtime.SetFinalizer(v, (*Vocoder).Close)

	// 初始化
	if ret := C.md380_vocoder_init(handle); ret != 0 {
		errMsg := v.getError()
		C.md380_vocoder_destroy(handle)
		return nil, fmt.Errorf("failed to initialize vocoder: %s", errMsg)
	}

	return v, nil
}

// Close 释放 vocoder 资源
func (v *Vocoder) Close() error {
	if v.handle != nil {
		C.md380_vocoder_destroy(v.handle)
		v.handle = nil
		runtime.SetFinalizer(v, nil)
	}
	return nil
}

// getError 获取最后的错误信息
func (v *Vocoder) getError() string {
	if v.handle == nil {
		return "vocoder is closed"
	}

	cErr := C.md380_vocoder_get_error(v.handle)
	if cErr != nil {
		return C.GoString(cErr)
	}
	return "unknown error"
}

// Encode 将 PCM 编码为 AMBE
func (v *Vocoder) Encode(pcm []int16) ([]byte, error) {
	if len(pcm) != PCMFrameSize {
		return nil, fmt.Errorf("PCM frame must be %d samples, got %d", PCMFrameSize, len(pcm))
	}

	if v.handle == nil {
		return nil, errors.New("vocoder is closed")
	}

	ambe := make([]byte, AMBEFrameSize)

	ret := C.md380_vocoder_encode(
		v.handle,
		(*C.int16_t)(unsafe.Pointer(&pcm[0])),
		(*C.uint8_t)(unsafe.Pointer(&ambe[0])),
	)

	if ret != 0 {
		return nil, fmt.Errorf("encoding failed: %s", v.getError())
	}

	return ambe, nil
}

// Decode 将 AMBE 解码为 PCM
func (v *Vocoder) Decode(ambe []byte) ([]int16, error) {
	if len(ambe) != AMBEFrameSize {
		return nil, fmt.Errorf("AMBE frame must be %d bytes, got %d", AMBEFrameSize, len(ambe))
	}

	if v.handle == nil {
		return nil, errors.New("vocoder is closed")
	}

	pcm := make([]int16, PCMFrameSize)

	ret := C.md380_vocoder_decode(
		v.handle,
		(*C.uint8_t)(unsafe.Pointer(&ambe[0])),
		(*C.int16_t)(unsafe.Pointer(&pcm[0])),
	)

	if ret != 0 {
		return nil, fmt.Errorf("decoding failed: %s", v.getError())
	}

	return pcm, nil
}
