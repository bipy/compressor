English | [简体中文](README_zh_CN.md)

# compressor - ⚡️High-performance parallel image compression tool

> RAW images often have very large file sizes, taking up too much space in storage and organization.
>
> Image compression can significantly reduce image size without noticeably affecting image quality.

`compressor` implements the function of parallel compression for a large number of images, fully utilizing hardware performance and saving a lot of time.

![](https://goreportcard.com/badge/github.com/bipy/compressor)

# Features

- CLI
- High-performance parallel compression, with customizable parallel quantity
- Recursively accesses all images in the input folder
- Can specify output file path, or automatically generate in the parent directory of the image
- Automatically renames when encountering duplicate file names
- Supports adjusting output image quality
- Supports adjusting input format
- Supports adjusting output format
- Full logs
- Output statistics
- Exception handling
- Cross-platform support

# Usage

## Run

```bash
# CLI Mode
# 16 Threads; Quality 80; Input Path ~/Pictures
compressor -i ~/Pictures -j 16 -q 80

# Single File Mode
compressor -i ~/Pictures/test.png

# Use Webp
compressor -i ~/Pictures/test.png -t webp -q 75

# Help
compressor -h
```

## Example

**If the input folder is set to `~/Pictures/my-photos`, the program automatically generates the ID `1700457797`, and automatically creates the output path**

```
~/Pictures/my-photos -> ~/Pictures/my-photos-1700457797
~/Pictures/my-photos/part1 -> ~/Pictures/my-photos-1700457797/part1
```

**Recursively process all files**

```
~/Pictures/my-photos/part1/test.png -> ~/Pictures/my-photos-1700457797/part1/test.jpg
~/Pictures/my-photos/haha.png -> ~/Pictures/my-photos-1700457797/haha.jpg
```

**Automatic renaming**

```
~/Pictures/my-photos/haha.jpg -> ~/Pictures/my-photos-1700457797/haha-1.jpg
~/Pictures/my-photos/haha.jpeg -> ~/Pictures/my-photos-1700457797/haha-2.jpg
```

## Full Usage

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
