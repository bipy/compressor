package constant

const Version = "4.0"

const FormattedHelpMessage = `Version: %s
Usage: compressor [-h] [Options]

Options:
  -h            show this help
  -i <path>     input path (file or directory)
  -o <path>     output path (auto-generated if not specified)
  -t <format>   output format: jpg/jpeg/png/webp (default: jpg)
  -q <0-100>    output quality (default: 90)
  -j <number>   worker threads (default: CPU cores)
  -accept       accepted input formats (default: all)
  -width        max image width (default: original)
  -height       max image height (default: original)
`
