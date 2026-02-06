<div align="center">

# 🗜️ Compressor

### ⚡️ 高性能并行图片压缩工具

<p align="center">
  <strong>利用并行处理技术，在数秒内压缩数千张图片</strong>
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/bipy/compressor)](https://goreportcard.com/report/github.com/bipy/compressor)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://go.dev/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/bipy/compressor/pulls)

[English](README.md) | [简体中文](README_zh_CN.md)

---

</div>

## 📖 项目简介

RAW 格式图片通常占用大量存储空间，给归档和整理工作带来挑战。**Compressor** 通过利用并行处理技术来高效压缩大批量图片，在不明显影响图像质量的前提下大幅减少文件体积。

### ✨ 为什么选择 Compressor？

- 🚀 **极速处理** - 并行压缩充分利用 CPU 多核性能
- 📦 **批量处理** - 轻松处理数千张图片
- 🎯 **智能输出** - 自动生成有序的文件夹结构
- 🔄 **格式灵活** - 支持 JPG、PNG、WebP 格式转换
- 🛡️ **安全可靠** - 智能处理重名文件，具备错误恢复机制
- 🌍 **跨平台** - 完美支持 Linux、macOS 和 Windows

## 🎯 核心特性

| 功能 | 说明 |
|------|------|
| 🖥️ **命令行界面** | 简洁的命令行接口，便于自动化集成 |
| ⚡ **并行处理** | 可自定义线程数量，优化处理性能 |
| 🔍 **递归扫描** | 自动处理所有子目录中的图片 |
| 📂 **智能输出** | 自动生成有序的输出文件夹或使用自定义路径 |
| 🔄 **重名处理** | 遇到文件冲突时智能重命名 |
| 🎨 **质量控制** | 可调节压缩质量（0-100） |
| 📥 **格式筛选** | 指定要处理的输入格式 |
| 📤 **格式转换** | 转换为 JPG、PNG 或 WebP 格式 |
| 📊 **详细日志** | 全面的操作日志和统计信息 |
| 🛡️ **异常处理** | 健壮的异常处理机制，确保处理不中断 |
| 🌐 **跨平台** | 原生支持 Linux、macOS 和 Windows |

---

## 🚀 快速开始

### 安装方式

**使用 Go 安装：**
```bash
go install github.com/bipy/compressor@latest
```

**从源码编译：**
```bash
git clone https://github.com/bipy/compressor.git
cd compressor
go build -o compressor
```

### 基本使用

```bash
# 使用 16 线程，质量设为 80 进行压缩
compressor -i ~/Pictures -j 16 -q 80

# 压缩单个文件
compressor -i ~/Pictures/photo.png

# 转换为 WebP 格式
compressor -i ~/Pictures -t webp -q 75

# 显示帮助信息
compressor -h
```

---

## 📚 使用示例

### 示例 1：批量压缩

输入文件夹：`~/Pictures/my-photos`

程序自动生成唯一 ID 并创建有序的输出结构：

```
📁 ~/Pictures/my-photos          →  📁 ~/Pictures/my-photos-1700457797
  📁 part1                        →    📁 part1
    🖼️ test.png                   →      🖼️ test.jpg
  🖼️ haha.png                     →    🖼️ haha.jpg
```

### 示例 2：递归处理

自动处理所有子目录中的图片：

```
📁 输入                                          📁 输出
~/Pictures/my-photos/part1/test.png    →    ~/Pictures/my-photos-1700457797/part1/test.jpg
~/Pictures/my-photos/haha.png          →    ~/Pictures/my-photos-1700457797/haha.jpg
```

### 示例 3：智能重名处理

检测到重复文件名时，自动进行重命名：

```
~/Pictures/my-photos/haha.jpg     →    ~/Pictures/my-photos-1700457797/haha-1.jpg
~/Pictures/my-photos/haha.jpeg    →    ~/Pictures/my-photos-1700457797/haha-2.jpg
```

---

## ⚙️ 命令行选项

```
用法: compressor [-h] [选项]

选项:
  -h                    显示帮助信息
  -i <路径>             输入路径（文件或目录）
  -o <路径>             输出路径（可选，未指定时自动生成）
  -j <数字>             并行处理的线程数量（默认: 8）
  -q <0-100>           输出质量（默认: 90）
  -t <格式>             输出格式: jpg/jpeg/png/webp（默认: jpg）
  -accept <格式>        接受的输入格式（默认: "jpg jpeg png"）
  -width <像素>         最大图片宽度（默认: 原始尺寸）
  -height <像素>        最大图片高度（默认: 原始尺寸）
```

### 常用示例

```bash
# 使用最大线程数进行高质量压缩
compressor -i ~/Photos -j 32 -q 95

# 转换 PNG 为 WebP 并限制尺寸
compressor -i ~/Photos -t webp -width 1920 -height 1080

# 仅处理特定格式
compressor -i ~/Photos -accept "png jpg" -t jpg -q 85

# 指定自定义输出目录
compressor -i ~/Photos/raw -o ~/Photos/compressed -j 16
```

---

## 🏗️ 工作原理

1. **扫描** - 递归发现输入路径中的所有图片
2. **筛选** - 根据 `-accept` 参数应用格式过滤
3. **处理** - 使用工作线程并行压缩图片
4. **转换** - 转换为目标格式（JPG/PNG/WebP）
5. **保存** - 将压缩后的图片写入输出目录
6. **报告** - 显示压缩统计信息和结果

---

## 🤝 参与贡献

欢迎贡献代码！请随时提交 Pull Request。对于重大更改，请先开启 issue 讨论您想要改变的内容。

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

基于以下优秀项目构建：
- [disintegration/imaging](https://github.com/disintegration/imaging) - 图片处理
- [go-webpbin](https://github.com/nickalie/go-webpbin) - WebP 支持
- [charmbracelet/log](https://github.com/charmbracelet/log) - 美观的日志输出

---

<div align="center">

**用 ❤️ 制作，来自 Compressor 团队**

⭐ 给我们的 GitHub 仓库点个星 — 这对我们很有激励作用！

[报告问题](https://github.com/bipy/compressor/issues) · [功能建议](https://github.com/bipy/compressor/issues)

</div>
