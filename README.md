# This project is based on imageflow

High-performance image manipulation for web servers. Includes imageflow_server, imageflow_tool, and libimageflow.

项目地址：[imazen/imageflow](https://github.com/imazen/imageflow)

官方网站：[imageflow = libimageflow + imageflow-server](https://www.imageflow.io/)



# 基于 imageflow 的多线程图片批量压缩脚本

imageflow 预编译的 image_tool 可以实现对图片进行高质量的压缩，可以在不显著影响图像质量的情况下大幅减少图片体积

但在大量处理图片的情形下需要脚本来完成并行与自动化处理等功能，因此选择了 Java 来完成并行和递归处理功能，与 libimageflow 的交互通过命令行工具 imageflow_toool 实现



# Feature

- 使用线程池并行处理，可自定义线程数
- 递归访问输入文件夹下所有图片
- 可以指定输出文件路径，也可以在图片父目录下自动生成
- 支持 imageflow 的主要参数，可用于调整编码格式，输出图片质量和格式，图片大小等
- 输出统计与异常处理



# Demo

作品名称：しのぶ 🦋

作品地址：[Pixiv-77458895](https://www.pixiv.net/artworks/77458895)

imageflow 版本：1.3.6-rc36

输出格式：JPG

输出质量参数：90

体积减小：85.3% (3,322 KB --> 487 KB)



## Original

Solution: 3000 × 1688 (PNG, RGBA32)

Size: 3.24 MB

![](https://cdn.jsdelivr.net/gh/bipy/CDN@master/repo/Image-Compressor/pid-77458895.png)



## Compressed

Solution: 3000 × 1688 (JPG, YUV)

Size: 486 KB

![](https://cdn.jsdelivr.net/gh/bipy/CDN@master/repo/Image-Compressor/pid-77458895_compressed.jpg)



# Usage

注意：Windows下不支持中文路径

**解决方法：**

1. 修改为英文路径
2. 控制面板 --> 时钟和区域 --> 区域 --> 管理 --> 更改系统区域设置 --> 使用Unicode UTF-8提供全球语言支持



## 环境配置

支持的 JDK 版本为 1.8+

需下载 [imageflow Releases](https://github.com/imazen/imageflow/releases) 并配置 imageflow_tool 的路径



## 默认的参数配置

将图片压缩为 jpg 格式

输出质量90%，不修改图片大小

输出路径为同级目录下的`compressed`文件夹

输出图片文件名后添加`_compressed`后缀



## 修改参数

直接修改`Variables.java`文件

```java
// 例如：

// 处理软件IMAGE FLOW位置
public static final String IMAGE_FLOW_TOOL_PATH = "D:\\imageflow_tool.exe";

// 输入路径，可以是文件夹（将递归处理所有子文件夹），如果没有全局utf-8的话要避免中文路径
public static final String INPUT_PATH = "D:\\Users\\Fidelity\\Pictures\\2020";

```



## 启动

```bash
# 进入src/ 确认参数后编译
javac -encoding UTF-8 Main.java
# 编译完成生成 class 文件

# 运行
java Main
```



# 变量说明

|         名称         |       可选值        |                             说明                             |
| :------------------: | :-----------------: | :----------------------------------------------------------: |
| IMAGE_FLOW_TOOL_PATH |    imageflow.exe    |                     imageflow_tool 路径                      |
|     THREAD_COUNT     | 小于CPU核心数的两倍 |                            线程数                            |
|      INPUT_PATH      |      D:\\in\\       | 输入路径，可以是文件夹（将递归处理所有子文件夹），如果没有全局utf-8的话要避免中文路径 |
|     OUTPUT_PATH      |      D:\\out\\      |     输出路径，默认为图片当前文件夹下新建compressed文件夹     |
|   AUTO_OUTPUT_PATH   |     true/false      | 自动选择输出路径，true: 默认设置；false: 在OUTPUT_PATH中设置指定文件夹 |
|   OUTPUT_PATH_NAME   |       String        |   自动生成的输出文件夹名称，AUTO_OUTPUT_PATH为false时无效    |
|     PROCESS_TYPE     |      v0.1/ir4       |                   运行模式，固定参数勿修改                   |
|       QUALITY        |       0 ~ 100       |                     输出图片质量，推荐90                     |
|    OUTPUT_FORMAT     |       jpg/png       |                       输出图片编码格式                       |
|        RESIZE        |     true/false      |                       是否改变图片大小                       |
|     FIXED_WIDTH      |     true/false      | true: 按指定宽度缩放，false: 按指定高度缩放，只会缩小，不会放大，RESIZE为false时无效 |
|        WIDTH         |         int         |                指定宽度，RESIZE为false时无效                 |
|        HEIGHT        |         int         |                指定高度，RESIZE为false时无效                 |
|  OUTPUT_PIC_POSTFIX  |       String        |                         输出图片后缀                         |



