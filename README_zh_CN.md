[English](README.md) | 简体中文

# compressor - ⚡️高性能并行图片压缩工具

> RAW 图片往往具有非常大的文件体积，在归纳整理存储上占据了太多的空间
>
> 图片压缩可以在不显著影响图像质量的情况下大幅减少图片体积

`compressor` 实现了对超大量图片并行压缩的功能，可充分利用硬件性能，节省大量时间

![](https://goreportcard.com/badge/github.com/bipy/compressor)

# 功能

- CLI
- 高性能并行压缩，可自定义并行数量
- 递归访问输入文件夹下所有图片
- 可以指定输出文件路径，也可以在图片父目录下自动生成
- 遇到重复文件名时自动重命名
- 支持调整输出图片质量
- 支持调整输入格式
- 支持调整输出格式
- 完整日志
- 输出统计
- 异常处理
- 跨平台支持

# 使用方法

## 启动

```bash
# CLI 模式
# 16 线程; 质量 80; 输入路径 ~/Pictures
compressor -i ~/Pictures -j 16 -q 80

# 单文件模式
compressor -i ~/Pictures/test.png

# 使用 Webp
compressor -i ~/Pictures/test.png -t webp -q 75

# 完整用法
compressor -h
```

## 示例

**若配置输入文件夹为`~/Pictures/my-photos`，程序自动生成ID `1700457797`，自动创建输出路径**

```
~/Pictures/my-photos -> ~/Pictures/my-photos-1700457797
~/Pictures/my-photos/part1 -> ~/Pictures/my-photos-1700457797/part1
```

**递归处理所有文件**

```
~/Pictures/my-photos/part1/test.png -> ~/Pictures/my-photos-1700457797/part1/test.jpg
~/Pictures/my-photos/haha.png -> ~/Pictures/my-photos-1700457797/haha.jpg
```

**自动重命名**

```
~/Pictures/my-photos/haha.jpg -> ~/Pictures/my-photos-1700457797/haha-1.jpg
~/Pictures/my-photos/haha.jpeg -> ~/Pictures/my-photos-1700457797/haha-2.jpg
```

## 完整用法

```
Usage: compressor [-h] [Options]

Options:
  -h
        show this help
  -accept string
        accepted input format (default "jpg jpeg png")
  -height int
        max image height (default 9223372036854775807)
  -i string
        input path
  -j int
        thread count (default 8)
  -o string
        output path
  -q int
        output quality: 0-100 (default 90)
  -t string
        output type: jpg/jpeg/png/webp (default "jpg")
  -width int
        max image width (default 9223372036854775807)
```
