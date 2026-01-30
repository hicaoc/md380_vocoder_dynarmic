# MD380 Vocoder (CGo 封装)

本项目将 MD380 固件的声码器功能（基于 Dynarmic 模拟器）封装为 Go 语言库。项目包含一个预编译的 Linux x86_64 静态库，使得 Go 程序可以直接通过 CGo 链接使用，而无需在使用者机器上配置复杂的 C++ 编译环境。

## CGo 架构说明

核心的声码器逻辑由 C++ 实现，并编译为静态库 `libmd380_vocoder.a`。Go 封装层直接链接该库。

- **Go 包名**: `github.com/hicaoc/md380_vocoder_dynarmic`
- **链接配置**: 包内自动配置了 CGo 指令，链接位于 `go_wrapper/lib` 目录下的静态库。
  ```go
  // #cgo LDFLAGS: -L${SRCDIR}/go_wrapper/lib -lmd380_vocoder -lstdc++ -lm
  // #cgo linux LDFLAGS: -ldl -lpthread
  ```
- **依赖管理**: 静态库已将所有必要的依赖项（`dynarmic`, `fmt`, `mcl`, `Zydis`, `Zycore`）打包在内，确保使用者可以“开箱即用”，无需手动安装这些 C++ 库。

## 验证与测试报告 (Test Report)

本封装库已通过单元测试验证，确保编码/解码功能正确，且内存操作安全。

### 测试环境
- **操作系统**: Linux x86_64
- **Go 版本**: 1.25.6
- **测试日期**: 2026-01-30

### 测试结果详情
```text
=== RUN   TestVocoder_Encode
md380_init: firmware=0xb19f00 sram=0xaf9d20
Encoded AMBE frame: aadc81128f8d0988ad
--- PASS: TestVocoder_Encode (0.09s)
    [通过] 成功将 160 个 PCM 采样编码为有效的 9 字节 AMBE 帧。

=== RUN   TestVocoder_EncodeDecode
md380_init: firmware=0xb19f00 sram=0xaf9d20
Decoded first sample: 0 (orig: 0)
--- PASS: TestVocoder_EncodeDecode (0.07s)
    [通过] 编码->解码 闭环测试通过，帧长保持一致 (160 采样 -> 9 字节 -> 160 采样)。

=== RUN   TestVocoder_EncodeSilence
md380_init: firmware=0xb19f00 sram=0xaf9d20
Encoded Silence AMBE: bbea81724272006a6f
--- PASS: TestVocoder_EncodeSilence (0.06s)
    [通过] 静音帧编码测试成功。

PASS
ok      github.com/hicaoc/md380_vocoder_dynarmic        0.223s
```

## 使用方法 (Usage)

```go
import "github.com/hicaoc/md380_vocoder_dynarmic"

func main() {
    // 1. 初始化声码器实例
    v, err := md380vocoder.NewVocoder()
    if err != nil {
        panic(err)
    }
    
    // 2. 编码: 输入 160 个 int16 采样
    pcm := make([]int16, 160)
    // 填充 pcm 数据...
    ambe, err := v.Encode(pcm)
    if err != nil {
        panic(err)
    }
    
    // 3. 解码: 输入 9 字节 AMBE 帧
    decodedPcm, err := v.Decode(ambe)
}
```

## 重新编译静态库 (高级选项)
如果你需要修改底层的 C++ 代码或重新生成 `libmd380_vocoder.a`，请按照以下步骤操作：

```bash
# 1. 创建构建目录
mkdir build
cd build

# 2. 配置 CMake (需要系统已安装 cmake, make, g++)
cmake ..

# 3. 编译完整静态库 (包含所有依赖)
make -j4 md380_vocoder_full

# 4. 更新 Go 封装库使用的静态文件
cp libmd380_vocoder_full.a ../go_wrapper/lib/libmd380_vocoder.a
```
*注意: 重新编译需要你的环境具备 `dynarmic` 的编译依赖 (fmt, mcl, Zydis 等)。普通使用者无需执行此步骤。*
