# 本项目基于 imageflow

High-performance image manipulation for web servers. Includes imageflow_server, imageflow_tool, and libimageflow.

项目地址：[imazen/imageflow](https://github.com/imazen/imageflow)

官方网站：[imageflow = libimageflow + imageflow-server](https://www.imageflow.io/)



# 基于 imageflow 的多线程图片批量压缩脚本

imageflow 预编译的 imageflow_tool 可以实现对图片进行高质量的压缩，可以在不显著影响图像质量的情况下大幅减少图片体积

但在大量处理图片的情形下需要脚本来完成并行与自动化处理等功能，因此选择了 Go 来完成并行和递归处理功能，与 libimageflow 的交互通过命令行工具 imageflow_tool 实现



# 功能

- 使用协程并行调用 imageflow_tool，可自定义线程数
- 递归访问输入文件夹下所有图片
- 可以指定输出文件路径，也可以在图片父目录下自动生成
- 支持 imageflow 的主要参数，可用于调整编码格式，输出图片质量和格式，图片大小等
- 完整日志保存
- 输出统计
- 异常处理
- 支持 Windows / Linux / Mac OS 等



# DEMO

作品名称：しのぶ 🦋

作品地址：[Pixiv-77458895](https://www.pixiv.net/artworks/77458895)

imageflow 版本：1.5.0-rc54

输出格式：JPG

输出质量参数：90

体积减小：85.9% (3,322 KB -> 468 KB)



## 原始

格式: 3000 × 1688 (PNG, RGBA32)

文件体积: 3.24 MB

![](https://cdn.jsdelivr.net/gh/bipy/CDN@master/repo/Image-Compressor/pid-77458895.png)



## 压缩后

格式: 3000 × 1688 (JPG, YUV420)

文件体积: 468 KB

![](https://cdn.jsdelivr.net/gh/bipy/CDN@master/repo/Image-Compressor/pid-77458895.jpg)



# 用法

注意：Windows 下不支持中文路径 (**经测试，在最新版本 imageflow 中已支持中文**)

**解决方法：**

1. 修改为英文路径
2. 控制面板 --> 时钟和区域 --> 区域 --> 管理 --> 更改系统区域设置 --> 使用 Unicode UTF-8 提供全球语言支持



## imageflow 已知 BUG

- 不能处理宽或高大于 10000 像素的图片
- v1.5.5-rc59 存在处理极慢（耗时几分钟）的情况，请选择其他版本



## 配置

需下载 [imageflow Releases](https://github.com/imazen/imageflow/releases) 并配置 imageflow_tool 的路径



## 默认的参数配置

将图片压缩为 **jpg** 格式

输出质量 **90%**

**不修改**图片大小

自动生成 ID 以区分

自动输出路径为父级目录下的 `INPUT_ID` 文件夹

保留所有图片的相对路径

**eg：**若配置输入文件夹为`D:\\Pictures`

```
D:\\Pictures\\myimg\\test.png -> D:\\Pictures_231453823\\myimg\\test.jpg
D:\\Pictures\\mypic\\hahaha.png -> D:\\Pictures_231453823\\mypic\\hahaha.jpg
```



## 修改参数

直接修改 `config.json` 文件



## 启动

下载对应 Release 并配置

```bash
# 无参数，使用相对路径下的 config.json
compressor

# 指定配置文件
compressor -c another_config.json
```



# 配置文件说明

|      名称      |       可选值        |                         说明                         |
| :------------: | :-----------------: | :--------------------------------------------------: |
| imageflow_tool |    imageflow.exe    |                 imageflow_tool 路径                  |
|  thread_count  | 小于CPU核心数的两倍 |                        线程数                        |
|   input_path   |      D:\\in\\       |   输入路径，必须是文件夹（将递归处理所有子文件夹）   |
|  output_path   |      D:\\out\\      | 输出路径，若不指定默认在输入路径父文件夹下新建文件夹 |
|  output_image  |         {}          |                   图片输出相关设置                   |

## 图片输出相关设置（output_image）

|     名称      | 可选值  |         说明          |
| :-----------: | :-----: | :-------------------: |
|    quality    | 0\~100  | 输出图片质量，推荐 90 |
| output_format | jpg/png |   输出图片编码格式    |
|    resize     |   {}    |   图片尺寸相关设置    |

## 图片尺寸相关设置（resize）

|   名称    |   可选值   |                 说明                 |
| :-------: | :--------: | :----------------------------------: |
|  enable   | true/false |           是否改变图片尺寸           |
| resize_by |    0/1     | 0: 按指定宽度缩放，1: 按指定高度缩放 |
|   width   |    int     |             指定图片宽度             |
|  height   |    int     |             指定图片高度             |