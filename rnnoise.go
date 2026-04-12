package rnnoise

/*
#include <rnnoise.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

var (
	ErrAlreadyInitialized = errors.New("already initialized")
	ErrNotInitialized     = errors.New("not initialized")
	ErrBufTooSmall        = errors.New("buffer must be at least rnnoise.GetFrameSize() large")
)

func GetFrameSize() int {
	return int(C.rnnoise_get_frame_size())
}

type DenoiseState struct {
	p *C.struct_DenoiseState
	// allow GC to clean up the DenoiseState
	mem []byte
}

func NewDenoiseState() (*DenoiseState, error) {
	var d DenoiseState
	err := d.Init()
	return &d, err
}

func (d *DenoiseState) Init() error {
	if d.p != nil {
		return ErrAlreadyInitialized
	}
	size := C.rnnoise_get_size()
	d.mem = make([]byte, size)
	d.p = (*C.DenoiseState)(unsafe.Pointer(&d.mem))
	C.rnnoise_init(d.p, nil)
	return nil
}

// Decodes a frame of samples
// in and out must be at least rnnoise.GetFrameSize() large
func (d *DenoiseState) ProcessFrame(out []float32, in []float32) error {
	if d.p == nil {
		return ErrNotInitialized
	}
	fs := GetFrameSize()
	if len(out) < fs || len(in) < fs {
		return ErrBufTooSmall
	}

	outp := (*C.float)(unsafe.Pointer(&out))
	inp := (*C.float)(unsafe.Pointer(&in))
	C.rnnoise_process_frame(d.p, outp, inp)
	return nil
}
