// Package audio captures microphone input via miniaudio (malgo) at the 16 kHz
// mono format whisper.cpp expects, and encodes it as WAV for the STT sidecar.
package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/gen2brain/malgo"
)

const (
	SampleRate    = 16000
	Channels      = 1
	BitsPerSample = 16 // malgo.FormatS16
)

// Recorder is an in-progress capture. Create with Start, finish with Stop.
type Recorder struct {
	ctx    *malgo.AllocatedContext
	device *malgo.Device

	mu  sync.Mutex
	pcm []byte
}

// Start opens the default input device and begins capturing. miniaudio resamples
// the hardware stream down to 16 kHz mono for us.
func Start() (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("audio context init failed: %w", err)
	}

	r := &Recorder{ctx: ctx, pcm: make([]byte, 0, SampleRate*2*4)}

	cfg := malgo.DefaultDeviceConfig(malgo.Capture)
	cfg.Capture.Format = malgo.FormatS16
	cfg.Capture.Channels = Channels
	cfg.SampleRate = SampleRate
	cfg.Alsa.NoMMap = 1

	callbacks := malgo.DeviceCallbacks{
		Data: func(_, input []byte, _ uint32) {
			r.mu.Lock()
			r.pcm = append(r.pcm, input...)
			r.mu.Unlock()
		},
	}

	device, err := malgo.InitDevice(ctx.Context, cfg, callbacks)
	if err != nil {
		ctx.Uninit()
		ctx.Free()
		return nil, fmt.Errorf("audio device init failed: %w", err)
	}
	r.device = device

	if err := device.Start(); err != nil {
		device.Uninit()
		ctx.Uninit()
		ctx.Free()
		return nil, fmt.Errorf("failed to start capture: %w", err)
	}
	return r, nil
}

// Stop ends capture and returns the audio encoded as a 16 kHz mono WAV.
func (r *Recorder) Stop() []byte {
	if r.device != nil {
		r.device.Uninit() // stops and frees
	}
	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
	}
	r.mu.Lock()
	pcm := r.pcm
	r.mu.Unlock()
	return EncodeWAV(pcm, SampleRate, Channels, BitsPerSample)
}

// DurationSeconds reports how much audio has been captured so far.
func (r *Recorder) DurationSeconds() float64 {
	r.mu.Lock()
	n := len(r.pcm)
	r.mu.Unlock()
	return float64(n) / float64(SampleRate*Channels*BitsPerSample/8)
}

// EncodeWAV wraps raw little-endian PCM in a canonical 44-byte WAV header.
func EncodeWAV(pcm []byte, sampleRate, channels, bitsPerSample int) []byte {
	var buf bytes.Buffer
	dataLen := len(pcm)
	byteRate := sampleRate * channels * bitsPerSample / 8
	blockAlign := channels * bitsPerSample / 8

	buf.WriteString("RIFF")
	binary.Write(&buf, binary.LittleEndian, uint32(36+dataLen))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	binary.Write(&buf, binary.LittleEndian, uint32(16)) // PCM fmt chunk size
	binary.Write(&buf, binary.LittleEndian, uint16(1))  // audio format = PCM
	binary.Write(&buf, binary.LittleEndian, uint16(channels))
	binary.Write(&buf, binary.LittleEndian, uint32(sampleRate))
	binary.Write(&buf, binary.LittleEndian, uint32(byteRate))
	binary.Write(&buf, binary.LittleEndian, uint16(blockAlign))
	binary.Write(&buf, binary.LittleEndian, uint16(bitsPerSample))

	buf.WriteString("data")
	binary.Write(&buf, binary.LittleEndian, uint32(dataLen))
	buf.Write(pcm)

	return buf.Bytes()
}
