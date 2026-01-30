package md380vocoder

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"
)

func TestVocoder_Encode(t *testing.T) {
	v, err := NewVocoder()
	if err != nil {
		t.Fatalf("Failed to create vocoder: %v", err)
	}
	// Verify Close doesn't panic
	defer func() {
		if err := v.Close(); err != nil {
			t.Errorf("Failed to close vocoder: %v", err)
		}
	}()

	// Generate a 1kHz sine wave
	// Sample rate is typically 8000Hz for these vocoders, but frame size 160 implies 20ms.
	pcm := make([]int16, PCMFrameSize)
	for i := 0; i < len(pcm); i++ {
		// Just some signal
		pcm[i] = int16(3000 * math.Sin(2*math.Pi*float64(i)*1000.0/8000.0))
	}

	ambe, err := v.Encode(pcm)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(ambe) != AMBEFrameSize {
		t.Errorf("Expected AMBE frame size %d, got %d", AMBEFrameSize, len(ambe))
	}

	fmt.Printf("Encoded AMBE frame: %s\n", hex.EncodeToString(ambe))
}

func TestVocoder_EncodeDecode(t *testing.T) {
	v, err := NewVocoder()
	if err != nil {
		t.Fatalf("Failed to create vocoder: %v", err)
	}
	defer v.Close()

	// 1. PCM -> Encode -> AMBE
	pcmOrig := make([]int16, PCMFrameSize)
	for i := 0; i < len(pcmOrig); i++ {
		pcmOrig[i] = int16(1000 * math.Sin(2*math.Pi*float64(i)*400.0/8000.0))
	}

	ambe, err := v.Encode(pcmOrig)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// 2. AMBE -> Decode -> PCM
	pcmDecoded, err := v.Decode(ambe)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(pcmDecoded) != PCMFrameSize {
		t.Errorf("Expected decoded PCM size %d, got %d", PCMFrameSize, len(pcmDecoded))
	}

	// Note: Lossy codec means pcmOrig != pcmDecoded.
	// We can't assert equality. We just check if it runs.
	fmt.Printf("Decoded first sample: %d (orig: %d)\n", pcmDecoded[0], pcmOrig[0])
}

func TestVocoder_EncodeSilence(t *testing.T) {
	v, err := NewVocoder()
	if err != nil {
		t.Fatalf("Failed to create vocoder: %v", err)
	}
	defer v.Close()

	pcm := make([]int16, PCMFrameSize) // All zeros
	ambe, err := v.Encode(pcm)
	if err != nil {
		t.Fatalf("Encode silence failed: %v", err)
	}

	fmt.Printf("Encoded Silence AMBE: %s\n", hex.EncodeToString(ambe))
}
