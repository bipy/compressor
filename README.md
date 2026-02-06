<div align="center">

# 🗜️ Compressor

### ⚡️ High-Performance Parallel Image Compression Tool

<p align="center">
  <strong>Compress thousands of images in seconds with blazing-fast parallel processing</strong>
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/bipy/compressor)](https://goreportcard.com/report/github.com/bipy/compressor)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://go.dev/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/bipy/compressor/pulls)

[English](README.md) | [简体中文](README_zh_CN.md)

---

</div>

## 📖 Overview

RAW images often consume massive amounts of storage space, making organization and archiving challenging. **Compressor** solves this problem by leveraging parallel processing to compress large batches of images efficiently, reducing file sizes dramatically without noticeable quality loss.

### ✨ Why Compressor?

- 🚀 **Blazing Fast** - Parallel processing fully utilizes your CPU cores
- 📦 **Batch Processing** - Handle thousands of images effortlessly  
- 🎯 **Smart Output** - Auto-generates organized folder structures
- 🔄 **Format Flexible** - Supports JPG, PNG, WebP conversion
- 🛡️ **Safe & Reliable** - Smart duplicate handling and error recovery
- 🌍 **Cross-Platform** - Works on Linux, macOS, and Windows

## 🎯 Features

| Feature | Description |
|---------|-------------|
| 🖥️ **CLI Interface** | Simple command-line interface for easy automation |
| ⚡ **Parallel Processing** | Customizable thread count for optimal performance |
| 🔍 **Recursive Scanning** | Automatically processes all images in subdirectories |
| 📂 **Smart Output** | Auto-generates organized output folders or use custom paths |
| 🔄 **Duplicate Handling** | Intelligent renaming when file conflicts occur |
| 🎨 **Quality Control** | Adjustable compression quality (0-100) |
| 📥 **Format Filtering** | Specify which input formats to process |
| 📤 **Format Conversion** | Convert to JPG, PNG, or WebP |
| 📊 **Detailed Logging** | Comprehensive operation logs and statistics |
| 🛡️ **Error Handling** | Robust exception handling for uninterrupted processing |
| 🌐 **Cross-Platform** | Native support for Linux, macOS, and Windows |

---

## 🚀 Quick Start

### Installation

**Using Go:**
```bash
go install github.com/bipy/compressor@latest
```

**From Source:**
```bash
git clone https://github.com/bipy/compressor.git
cd compressor
go build -o compressor
```

### Basic Usage

```bash
# Compress images with 16 threads at quality 80
compressor -i ~/Pictures -j 16 -q 80

# Compress a single file
compressor -i ~/Pictures/photo.png

# Convert to WebP format
compressor -i ~/Pictures -t webp -q 75

# Show help
compressor -h
```

---

## 📚 Usage Examples

### Example 1: Batch Compression

Input folder: `~/Pictures/my-photos`

The program automatically generates a unique ID and creates an organized output structure:

```
📁 ~/Pictures/my-photos         →    📁 ~/Pictures/my-photos-1700457797
📁 part1                        →    📁 part1
🖼️ test.png                     →    🖼️ test.jpg
🖼️ haha.png                     →    🖼️ haha.jpg
```

### Example 2: Recursive Processing

All images in subdirectories are processed automatically:

```
📁 Input                                    📁 Output
~/Pictures/my-photos/part1/test.png    →    ~/Pictures/my-photos-1700457797/part1/test.jpg
~/Pictures/my-photos/haha.png          →    ~/Pictures/my-photos-1700457797/haha.jpg
```

### Example 3: Smart Duplicate Handling

When duplicate filenames are detected, automatic renaming kicks in:

```
~/Pictures/my-photos/haha.jpg     →    ~/Pictures/my-photos-1700457797/haha-1.jpg
~/Pictures/my-photos/haha.jpeg    →    ~/Pictures/my-photos-1700457797/haha-2.jpg
```

---

## ⚙️ Command Line Options

```
Usage: compressor [-h] [Options]

Options:
  -h                    Show this help message
  -i <path>            Input path (file or directory)
  -o <path>            Output path (optional, auto-generated if not specified)
  -j <number>          Thread count for parallel processing (default: 8)
  -q <0-100>           Output quality (default: 90)
  -t <format>          Output format: jpg/jpeg/png/webp (default: jpg)
  -accept <formats>    Accepted input formats (default: "jpg jpeg png")
  -width <pixels>      Maximum image width (default: original size)
  -height <pixels>     Maximum image height (default: original size)
```

### Common Examples

```bash
# High-quality compression with maximum threads
compressor -i ~/Photos -j 32 -q 95

# Convert PNG to WebP with size limit
compressor -i ~/Photos -t webp -width 1920 -height 1080

# Process only specific formats
compressor -i ~/Photos -accept "png jpg" -t jpg -q 85

# Specify custom output directory
compressor -i ~/Photos/raw -o ~/Photos/compressed -j 16
```

---

## 🏗️ How It Works

1. **Scan** - Recursively discovers all images in the input path
2. **Filter** - Applies format filters based on `-accept` parameter
3. **Process** - Compresses images in parallel using worker threads
4. **Convert** - Transforms to target format (JPG/PNG/WebP)
5. **Save** - Writes compressed images to output directory
6. **Report** - Displays compression statistics and results

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

Built with:
- [disintegration/imaging](https://github.com/disintegration/imaging) - Image processing
- [go-webpbin](https://github.com/nickalie/go-webpbin) - WebP support
- [charmbracelet/log](https://github.com/charmbracelet/log) - Beautiful logging

---

<div align="center">

**Made with ❤️**

⭐ Star on GitHub — it motivates a lot!

[Report Bug](https://github.com/bipy/compressor/issues) · [Request Feature](https://github.com/bipy/compressor/issues)

</div>
