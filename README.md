English | [简体中文](README_zh_CN.md)

# This project is based on imageflow

High-performance image manipulation for web servers. Includes imageflow_server, imageflow_tool, and libimageflow.

Github：[imazen/imageflow](https://github.com/imazen/imageflow)

Website：[imageflow = libimageflow + imageflow-server](https://www.imageflow.io/)



# Multi-thread image batch compression script based on imageflow

Imageflow's pre-compiled **imageflow_tool** can achieve high-quality image compression, which can greatly reduce the file size without significantly affecting the image quality.

However, in the case of processing a large number of images, scripts are needed to complete the functions of parallel and automatic processing, so Golang is selected to perform parallel and recursive processing functions.

Interaction with libimageflow is done through the command line tool imageflow_tool.



# Features

- Call imageflow_tool in parallel using goroutines, with a customizable number of threads
- Recursive access to all images under the input folder
- Output file paths can be specified, or generated automatically in the parent directory of the input folder
- Support for imageflow's main parameters, which can be used to adjust the encoding format, output image quality and format, image size, etc
- Full log
- Output statistics
- Exception handling
- Support Windows / Linux / Mac OS, etc



# DEMO

Title：しのぶ 🦋

Link：[Pixiv-77458895](https://www.pixiv.net/artworks/77458895)

imageflow Version：1.5.0-rc54

Output Format：JPG

Output Quality：90

Reduction：85.9% (3,322 KB -> 468 KB)



## Original

Solution: 3000 × 1688 (PNG, RGBA32)

Size: 3.24 MB

![](https://cdn.jsdelivr.net/gh/bipy/CDN@master/repo/Image-Compressor/pid-77458895.png)



## Compressed

Solution: 3000 × 1688 (JPG, YUV420)

Size: 468 KB

![](https://cdn.jsdelivr.net/gh/bipy/CDN@master/repo/Image-Compressor/pid-77458895.jpg)



# Usage



## Known BUGs about imageflow

- Cannot handle images with a width or height greater than 10000 pixels.
- v1.5.5-rc59 is extremely slow (it takes a few minutes per image). Please select another version.



## **Configuration**

Download [imageflow Releases](https://github.com/imazen/imageflow/releases) and configure the path of imageflow_tool



## Default

Compress the picture into **jpg** format

Output quality **90%**

will **not** resize the image

Automatically generate ID to distinguish

The automatic output path is the `INPUT_ID` folder under the parent directory of the input folder

Keep the relative paths of all images

**eg：**If the input folder is `D:\\Pictures`

```
D:\\Pictures\\myimg\\test.png -> D:\\Pictures_231453823\\myimg\\test.jpg
D:\\Pictures\\mypic\\hahaha.png -> D:\\Pictures_231453823\\mypic\\hahaha.jpg
```



## Modify configuration

Modify the file `config.json` directly



## Run

Download the release and configure

```bash
# No argument, using config.json under the relative path
compressor

# Specify configuration
compressor -c another_config.json
```



# Configuration description

|      Name      |             Value             |                         Description                          |
| :------------: | :---------------------------: | :----------------------------------------------------------: |
| imageflow_tool |         imageflow.exe         |                     imageflow_tool path                      |
|  thread_count  | lower than the number of core |                         Thread count                         |
|   input_path   |           D:\\in\\            | Input path, must be a folder (will recursively process all subfolders) |
|  output_path   |           D:\\out\\           | Output path,  create a new folder under the parent directory of the input folder if not specified |
|  output_image  |              {}               |                    Image related settings                    |

## Image related settings (output_image)

|     Name      |  Value  |                     Description                     |
| :-----------: | :-----: | :-------------------------------------------------: |
|    quality    | 0\~100  | Determines the jpeg encoding quality. Default is 90 |
| output_format | jpg/png |    Determines the format to encode the image as.    |
|    resize     |   {}    |             Image size related settings             |

## Image size related settings（resize）

|   Name    |   Value    |                         Description                          |
| :-------: | :--------: | :----------------------------------------------------------: |
|  enable   | true/false |                        resize or not                         |
| resize_by |    0/1     | 0: constrains the image by width，1: constrains the image by height |
|   width   |    int     |                        specify width                         |
|  height   |    int     |                        specify height                        |