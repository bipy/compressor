English | [简体中文](README_zh_CN.md)

# Multi-thread Image Batch Compression Tool

**jpeg** can achieve high-quality image compression, which can greatly reduce the file size without significantly affecting the image quality.

However, in the case of processing a large number of pictures, it is necessary to complete the functions of parallel and automatic processing, so Go is selected to build a high-performance image batch processing tool.

Compressor can properly handle the task of compressing a large number of pictures into JPG format.

![](https://goreportcard.com/badge/github.com/bipy/compressor)

# Features

- Compress in parallel using goroutines, with a customizable number of threads
- Recursive access to all images under the input folder
- Output file paths can be specified, or generated automatically in the parent directory of the input folder
- The output image quality can be adjusted, the output format is fixed to `.jpg`.
- Full log
- Output statistics
- Exception handling
- Cross-platform support

# Usage

**Default** compress the picture into `jpg` format

Output quality **90%**

Will **not** resize the image

Automatically generate ID to distinguish

The automatic output path is the `INPUT_ID` folder under the parent directory of the input folder

Keep the relative paths of all images

**eg**: If the input folder is `D:\\Pictures`

```
D:\\Pictures\\myimg\\test.png -> D:\\Pictures_231453823\\myimg\\test.jpg
D:\\Pictures\\mypic\\hahaha.png -> D:\\Pictures_231453823\\mypic\\hahaha.jpg
```

## Run

Download the release and configure

```bash
# Specify configuration
compressor -c config.json

# CLI Mode
# 8 threads; quality 80; input path ~/Pictures
compressor -i ~/Pictures -j 8 -q 80

# full usage
compressor -h
```

## Full Usage

```
Version: 2.3
Usage: compressor [-h] [Options]

Options:
  -h
    	show this help
  -c string
    	Configuration Filepath
  -f string
    	Input Format (default "jpg jpeg png")
  -i string
    	Input Path
  -j int
    	Thread Count (default 4)
  -log
    	Save Log as File
  -o string
    	Output Path
  -q int
    	JPEG Quality (default 90)
```



## Modify Configuration

Modify the file `config.json` directly

# Configuration Description

|     Name     |             Value             |                                            Description                                            |
|:------------:|:-----------------------------:|:-------------------------------------------------------------------------------------------------:|
| thread_count | lower than the number of core |                                           Thread count                                            |
| input_format |     "png", "jpg", "jpeg"      |                        Input image format (other formats will be ignored)                         |
|  input_path  |           D:\\in\\            |          Input path, can be a file or a folder (will recursively process all subfolders)          |
| output_path  |           D:\\out\\           | Output path,  create a new folder under the parent directory of the input folder if not specified |
|   quality    |             1～100             |                        Determines the jpeg encoding quality. Default is 90                        |
| log_to_file  |             bool              |                                         save log to file                                          |


