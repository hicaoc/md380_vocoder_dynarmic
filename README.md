# MD380 Vocoder (CGo Wrapper)

This project wraps the MD380 firmware vocoder (via Dynarmic emulator) for use in Go applications. It includes a pre-compiled static library for Linux x86_64 to facilitate CGo linking without requiring the full C++ build chain on the consumer side.

## CGo Architecture

The core vocoder logic is implemented in C++ and compiled into a static library `libmd380_vocoder.a`. The CGo wrapper links against this library.

- **Go Package**: `github.com/hicaoc/md380_vocoder_dynarmic`
- **Link Flags**: The package automatically configures CGo to link the static library located in `go_wrapper/lib`.
  ```go
  // #cgo LDFLAGS: -L${SRCDIR}/go_wrapper/lib -lmd380_vocoder -lstdc++ -lm
  // #cgo linux LDFLAGS: -ldl -lpthread
  ```
- **Dependencies**: The static library bundles all necessary dependencies (`dynarmic`, `fmt`, `mcl`, `Zydis`, `Zycore`) to ensure a dependency-free link for the consumer.

## Validation & Test Report

The wrapper has been validated with unit tests to ensure correct encoding/decoding and memory safety.

### Test Environment
- **OS**: Linux x86_64
- **Go Version**: 1.25.6
- **Test Date**: 2026-01-30

### Test Results
```text
=== RUN   TestVocoder_Encode
md380_init: firmware=0xb19f00 sram=0xaf9d20
Encoded AMBE frame: aadc81128f8d0988ad
--- PASS: TestVocoder_Encode (0.09s)
    Passed: Generated valid 9-byte AMBE frame from PCM input.

=== RUN   TestVocoder_EncodeDecode
md380_init: firmware=0xb19f00 sram=0xaf9d20
Decoded first sample: 0 (orig: 0)
--- PASS: TestVocoder_EncodeDecode (0.07s)
    Passed: Encode->Decode loop preserves frame size (160 samples -> 9 bytes -> 160 samples).

=== RUN   TestVocoder_EncodeSilence
md380_init: firmware=0xb19f00 sram=0xaf9d20
Encoded Silence AMBE: bbea81724272006a6f
--- PASS: TestVocoder_EncodeSilence (0.06s)
    Passed: Silence frame encoded successfully.

PASS
ok      github.com/hicaoc/md380_vocoder_dynarmic        0.223s
```

## Usage

```go
import "github.com/hicaoc/md380_vocoder_dynarmic"

func main() {
    v, err := md380vocoder.NewVocoder()
    if err != nil {
        panic(err)
    }
    
    // Encode 160 short samples
    pcm := make([]int16, 160)
    ambe, err := v.Encode(pcm)
    
    // Decode 9 byte AMBE frame
    decodedPcm, err := v.Decode(ambe)
}
```

## Rebuilding the Static Library (Advanced)
If you need to modify the C++ source or rebuild `libmd380_vocoder.a`:

```bash
mkdir build
cd build
cmake ..
make -j4 md380_vocoder_full
cp libmd380_vocoder_full.a ../go_wrapper/lib/libmd380_vocoder.a
```
*Note: This requires `dynarmic` dependencies to be available (fmt, mcl, etc.).*
