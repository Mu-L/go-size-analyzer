# go-size-analyzer

[![Go Report Card](https://goreportcard.com/badge/github.com/Zxilly/go-size-analyzer)](https://goreportcard.com/report/github.com/Zxilly/go-size-analyzer)
[![Tests](https://github.com/Zxilly/go-size-analyzer/actions/workflows/built-tests.yml/badge.svg)](https://github.com/Zxilly/go-size-analyzer/actions/workflows/built-tests.yml)
[![Codecov](https://img.shields.io/codecov/c/gh/Zxilly/go-size-analyzer)](https://codecov.io/github/Zxilly/go-size-analyzer)
[![GitHub release](https://img.shields.io/github/v/release/Zxilly/go-size-analyzer)](https://github.com/Zxilly/go-size-analyzer/releases)
[![Featured｜HelloGitHub](https://abroad.hellogithub.com/v1/widgets/recommend.svg?rid=c31575032ae144c6b22c553399a42742&claim_uid=QksOfglJ0njwLKa&theme=small)](https://hellogithub.com/repository/c31575032ae144c6b22c553399a42742)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/Zxilly/go-size-analyzer/badge)](https://scorecard.dev/viewer/?uri=github.com/Zxilly/go-size-analyzer)

一个简单的工具，用于分析 Go 编译二进制文件的大小。

- [x] 支持跨平台分析 `ELF`、`Mach-O`、`PE` 和 `WebAssembly(实验性)` 二进制格式
- [x] 按包和区段提供详细的大小分析
- [x] 支持多种输出格式: `text`、`json`、`html`、`svg`
- [x] 通过网页界面和终端 UI 进行交互式探索
- [x] 比较二进制的 diff 模式 (支持 `json` 和 `text` 输出)

## 安装

[![Packaging status](https://repology.org/badge/vertical-allrepos/go-size-analyzer.svg)](https://repology.org/project/go-size-analyzer/versions)

### [下载最新二进制文件](https://github.com/Zxilly/go-size-analyzer/releases)

### MacOS / Linux 通过 Homebrew 安装：

使用 [Homebrew](https://brew.sh/)
```
brew install go-size-analyzer
```

### Windows：

使用 [scoop](https://scoop.sh/)
```
scoop install go-size-analyzer
```

### Go 安装：
```
go install github.com/Zxilly/go-size-analyzer/cmd/gsa@latest
```

## 使用

### Example

#### Web mode

```bash
$ gsa --web golang-compiled-binary
```

将在 8080 端口启动一个 web 服务器，您可以在浏览器中查看结果。

或者您可以在浏览器中使用 WASM 版本：[GSA Treemap](https://gsa.zxilly.dev)

> [!NOTE]
> 由于浏览器的限制，wasm 版本比原生版本慢得多。
> 通常分析相同二进制文件需要 10 倍的时间。
> 
> 仅建议用于分析小型应用程序（大小小于 30 MB）

网页将如下所示：

![image](https://github.com/Zxilly/go-size-analyzer/assets/31370133/e69583ce-b189-4a0d-b108-c3b7d5c33a82)

您可以点击以展开包以查看详细信息。

#### 终端用户界面

```bash
$ gsa --tui golang-compiled-binary
```

[![asciicast](https://asciinema.org/a/670664.svg)](https://asciinema.org/a/670664)

#### 文本模式

```bash
$ gsa docker-compose-linux-x86_64
┌─────────────────────────────────────────────────────────────────────────────────┐
│ docker-compose-linux-x86_64                                                     │
├─────────┬──────────────────────────────────────────────────┬────────┬───────────┤
│ PERCENT │ NAME                                             │ SIZE   │ TYPE      │
├─────────┼──────────────────────────────────────────────────┼────────┼───────────┤
│ 17.37%  │ k8s.io/api                                       │ 11 MB  │ vendor    │
│ 15.52%  │ .rodata                                          │ 9.8 MB │ section   │
│ 8.92%   │ .gopclntab                                       │ 5.6 MB │ section   │
│ 7.51%   │ .strtab                                          │ 4.7 MB │ section   │
│ 5.13%   │ k8s.io/client-go                                 │ 3.2 MB │ vendor    │
│ 3.36%   │ .symtab                                          │ 2.1 MB │ section   │
│ 3.29%   │ github.com/moby/buildkit                         │ 2.1 MB │ vendor    │
│ 2.02%   │ google.golang.org/protobuf                       │ 1.3 MB │ vendor    │
│ 1.96%   │ github.com/google/gnostic-models                 │ 1.2 MB │ vendor    │
│ 1.82%   │ k8s.io/apimachinery                              │ 1.1 MB │ vendor    │
│ 1.73%   │ net                                              │ 1.1 MB │ std       │
│ 1.72%   │ github.com/aws/aws-sdk-go-v2                     │ 1.1 MB │ vendor    │
│ 1.57%   │ crypto                                           │ 991 kB │ std       │
│ 1.53%   │ github.com/docker/compose/v2                     │ 964 kB │ vendor    │
│ 1.48%   │ github.com/gogo/protobuf                         │ 931 kB │ vendor    │
│ 1.40%   │ runtime                                          │ 884 kB │ std       │
│ 1.32%   │ go.opentelemetry.io/otel                         │ 833 kB │ vendor    │
│ 1.28%   │ .text                                            │ 809 kB │ section   │
│ 1.18%   │ google.golang.org/grpc                           │ 742 kB │ vendor    │

...[Collapsed]...

│ 0.00%   │ github.com/google/shlex                          │ 0 B    │ vendor    │
│ 0.00%   │ github.com/pmezard/go-difflib                    │ 0 B    │ vendor    │
│ 0.00%   │ go.uber.org/mock                                 │ 0 B    │ vendor    │
│ 0.00%   │ github.com/kballard/go-shellquote                │ 0 B    │ vendor    │
│ 0.00%   │ tags.cncf.io/container-device-interface          │ 0 B    │ vendor    │
│ 0.00%   │ github.com/josharian/intern                      │ 0 B    │ vendor    │
│ 0.00%   │ github.com/shibumi/go-pathspec                   │ 0 B    │ vendor    │
│ 0.00%   │ dario.cat/mergo                                  │ 0 B    │ vendor    │
│ 0.00%   │ github.com/mattn/go-colorable                    │ 0 B    │ vendor    │
│ 0.00%   │ github.com/secure-systems-lab/go-securesystemslib│ 0 B    │ vendor    │
├─────────┼──────────────────────────────────────────────────┼────────┼───────────┤
│ 100%    │ KNOWN                                            │ 63 MB  │           │
│ 100%    │ TOTAL                                            │ 63 MB  │           │
└─────────┴──────────────────────────────────────────────────┴────────┴───────────┘

```

#### 差异模式

```bash
$ gsa bin-linux-1.21-amd64 bin-linux-1.22-amd64
┌────────────────────────────────────────────────────────────────┐
│ Diff between bin-linux-1.21-amd64 and bin-linux-1.22-amd64     │
├─────────┬──────────────────────┬──────────┬──────────┬─────────┤
│ PERCENT │ NAME                 │ OLD SIZE │ NEW SIZE │ DIFF    │
├─────────┼──────────────────────┼──────────┼──────────┼─────────┤
│ +28.69% │ runtime              │ 801 kB   │ 1.0 MB   │ +230 kB │
│ +100%   │ internal/chacha8rand │          │ 3.1 kB   │ +3.1 kB │
│ +5.70%  │ <autogenerated>      │ 18 kB    │ 19 kB    │ +1.0 kB │
│ +8.59%  │ internal/abi         │ 6.1 kB   │ 6.6 kB   │ +525 B  │
│ +10.52% │ internal/cpu         │ 4.9 kB   │ 5.4 kB   │ +515 B  │
│ +4.45%  │ internal/reflectlite │ 3.9 kB   │ 4.1 kB   │ +173 B  │
│ +2.64%  │ internal/bytealg     │ 1.5 kB   │ 1.5 kB   │ +39 B   │
│ +0.80%  │ strconv              │ 4.0 kB   │ 4.0 kB   │ +32 B   │
│ +0.19%  │ syscall              │ 13 kB    │ 13 kB    │ +24 B   │
│ -0.37%  │ embed                │ 8.6 kB   │ 8.6 kB   │ -32 B   │
│ -0.16%  │ main                 │ 19 kB    │ 19 kB    │ -32 B   │
│ -0.38%  │ reflect              │ 25 kB    │ 25 kB    │ -96 B   │
│ -0.26%  │ time                 │ 87 kB    │ 87 kB    │ -224 B  │
│ -7.95%  │ sync                 │ 9.5 kB   │ 8.7 kB   │ -755 B  │
├─────────┼──────────────────────┼──────────┼──────────┼─────────┤
│ +8.47%  │ .rodata              │ 122 kB   │ 132 kB   │ +10 kB  │
│ +5.04%  │ .gopclntab           │ 144 kB   │ 152 kB   │ +7.3 kB │
│ +3.61%  │ .debug_info          │ 168 kB   │ 174 kB   │ +6.1 kB │
│ +3.52%  │ .debug_loc           │ 81 kB    │ 84 kB    │ +2.9 kB │
│ +3.03%  │ .debug_line          │ 80 kB    │ 82 kB    │ +2.4 kB │
│ +3.41%  │ .symtab              │ 59 kB    │ 61 kB    │ +2.0 kB │
│ +4.29%  │ .debug_frame         │ 29 kB    │ 30 kB    │ +1.2 kB │
│ +1.25%  │ .strtab              │ 61 kB    │ 62 kB    │ +763 B  │
│ +3.28%  │ .debug_ranges        │ 13 kB    │ 13 kB    │ +415 B  │
│ +5.13%  │ .data                │ 5.0 kB   │ 5.2 kB   │ +256 B  │
│ +7.32%  │ .typelink            │ 1.3 kB   │ 1.3 kB   │ +92 B   │
│ +27.78% │ .go.buildinfo        │ 288 B    │ 368 B    │ +80 B   │
│ -1.56%  │ .debug_gdb_scripts   │ 64 B     │ 63 B     │ -1 B    │
│ -0.63%  │ .noptrdata           │ 2.5 kB   │ 2.5 kB   │ -16 B   │
├─────────┼──────────────────────┼──────────┼──────────┼─────────┤
│ +3.86%  │ bin-linux-1.21-amd64 │ 1.6 MB   │ 1.6 MB   │ +61 kB  │
│         │ bin-linux-1.22-amd64 │          │          │         │
└─────────┴──────────────────────┴──────────┴──────────┴─────────┘
```

#### Svg 模式

```bash
$ gsa cockroach-darwin-amd64 -f svg -o data.svg --hide-sections
```

![image](./assets/example.svg)

### 完整选项

```bash
A tool for determining the extent to which dependencies contribute to the
bloated size of compiled Go binaries.

Arguments:
  <file>           Binary file to analyze or result json file for diff
  [<diff file>]    New binary file or result json file to compare, optional

Flags:
  -h, --help             Show context-sensitive help.
      --verbose          Verbose output
  -f, --format="text"    Output format, possible values: text,json,html,svg
      --no-disasm        Skip disassembly pass
      --no-symbol        Skip symbol pass
      --no-dwarf         Skip dwarf pass
  -o, --output=STRING    Write to file
      --version          Show version

Text output options
  --hide-sections    Hide sections
  --hide-main        Hide main package
  --hide-std         Hide standard library

Json output options
  --indent=INDENT    Indentation for json output
  --compact          Hide function details, replacement with size

Svg output options
  --width=1028         Width of the svg treemap
  --height=640         Height of the svg treemap
  --margin-box=4       Margin between boxes
  --padding-box=4      Padding between box border and content
  --padding-root=32    Padding around root content

Web interface options
  --web               use web interface to explore the details
  --listen=":8080"    listen address
  --open              Open browser
  --update-cache      Update the cache file for the web UI

Terminal interface options
  --tui    Use terminal interface to explore the details

Imports analysis options
  --imports    Try analyze package imports from source

```

> [!CAUTION]
>
> 该工具可以分析剥离 symbol 的二进制文件，但可能导致结果不准确。

## TODO

- [ ] 添加更多用于反汇编二进制文件的模式
- [x] 从 DWARF 段提取信息
- [x] 计算符号本身的大小到包中
- [ ] 添加其他图表，如火焰图、饼图等
- [ ] 支持 demangle cgo 中的 C++/Rust 符号
- [x] 添加 TUI 模式以探索
- [x] 编译为 wasm，创建一个在浏览器中分析二进制文件的用户界面

## Contribution

欢迎任何形式的贡献，随时提出问题或拉取请求。

关于开发，请查看 [开发指南](./DEVELOPMENT.md) 以获取更多详细信息。

## LICENSE

根据 [AGPL-3.0](./LICENSE) 发布。
