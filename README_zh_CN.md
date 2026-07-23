<div align="center">

# 🗜️ Compressor

### 批量图片压缩，开箱即用

<br />

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-4caf50?style=flat-square)](LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/bipy/compressor?style=flat-square)](https://goreportcard.com/report/github.com/bipy/compressor)

[English](README.md) · [中文](README_zh_CN.md)

<br />

```
compressor -i ~/Photos -t webp -q 80
```

指向一个文件夹，剩下的交给它。

</div>

---

## 它能做什么

Compressor 拿到一堆图片，把它们变小。支持读取 JPEG、PNG、GIF、BMP、TIFF、WebP 和 HEIC，输出可以选择 JPEG、PNG 或 WebP。

它会并行处理，默认吃满你的 CPU。文件通过流式管道传输，扫描还没结束处理就已经开始了——不用干等。

## 安装

```bash
go install github.com/bipy/compressor@latest
```

<details>
<summary>想从源码编译？</summary>

```bash
git clone https://github.com/bipy/compressor.git
cd compressor
go build
```

</details>

## 试试看

```bash
# 整个文件夹
compressor -i ~/Photos

# 单个文件，转成 WebP
compressor -i photo.png -o photo.webp

# 缩放到 1920px 宽，质量 75
compressor -i ~/Photos -t webp -width 1920 -q 75

# 慢一点，对机器温柔些
compressor -i ~/Photos -j 2
```

## 所有参数

| 参数 | 默认值 | 作用 |
|:-----|:-------|:-----|
| `-i` | *（必填）* | 输入文件或目录 |
| `-o` | 自动生成 | 输出位置 |
| `-t` | `jpg` | 输出格式 — `jpg`、`png` 或 `webp` |
| `-q` | `90` | 质量，1 到 100 |
| `-j` | 全部 CPU 核心 | 工作线程数 |
| `-accept` | 全部格式 | 要处理哪些输入格式 |
| `-width` | 保持原始 | 最大宽度（像素） |
| `-height` | 保持原始 | 最大高度（像素） |

## 内部原理

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  发现    │───▶│  规划    │───▶│  执行    │───▶│  写入    │───▶│  报告    │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

**发现** — 遍历输入目录，把文件路径流式送进管道，不需要先完整扫描一遍。

**规划** — 检测每个文件的格式，算好输出放哪里，遇到重名自动处理。

**执行** — 把任务分发给工作池。每个 Worker 独立完成解码、缩放、编码、写入，互不干扰。

**写入** — 原子写入：先写临时文件，再重命名。中途失败也不会留下残缺文件。

**报告** — 汇总各项指标，实时推送到 TUI 仪表盘。如果是管道输出，自动切换纯文本。

## 终端界面

在终端里运行 Compressor，你会看到一个实时仪表盘：

```
Compressor
Total: 142 | ✓ 138 | ✗ 4 | Elapsed: 12s
Input: 2.1GB | Output: 856MB | Ratio: 0.41 | Speed: 178.5MB/s
[████████████████████░░░░░░░░░░] 67.2%

Recent:
  ✓ IMG_2024.jpg (2.1MB)
  ✓ IMG_2025.jpg (1.8MB)
  ✓ IMG_2026.jpg (2.3MB)

q: 退出 │ p: 暂停 │ v: 失败详情 │ ?: 帮助
```

管道输出？CI 环境？它会自动检测并切换到纯文本模式。

## 项目结构

```
internal/
├── codec/       编解码 — JPEG、PNG、WebP、GIF、BMP、TIFF、HEIC
├── discover/    流式文件遍历
├── executor/    工作池 + panic 恢复
├── model/       Job、Result、Summary 类型定义
├── planner/     格式检测、路径解析、冲突处理
├── reporter/    指标聚合
├── transform/   缩放和自动旋转
├── tui/         终端仪表盘（bubbletea）
└── writer/      原子文件写入
```

## 性能

Apple M4 基准测试（单线程，每张图片）：

| 操作 | 单张耗时 | 分配内存 | 含义 |
|:-----|--------:|---------:|:-----|
| 编码 JPEG | 271 μs | 7.2 KB | 约 3,700 张/秒 |
| 编码 PNG | 539 μs | 852 KB | 约 1,850 张/秒（PNG 更重） |
| 解码 JPEG | 152 μs | 83 KB | 约 6,500 张/秒 — 解码比编码快 |
| 扫描 1K 文件 | 33 μs | 7.6 KB | 1000 个文件扫描只需 33ms — 发现阶段几乎零开销 |

> **耗时** = 处理一张图片的时间，越小越好。
> **分配内存** = 每次操作的堆内存开销，越小 GC 压力越小。
> 以上是单 Worker 数据。实际吞吐量随 `-j` 线程数线性增长 — 8 核 ≈ 8 倍处理速度。

## 许可证

[MIT](LICENSE) — 随便用。
