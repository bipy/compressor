import java.io.File;

public class Picture {
    private File file;
    private String[] args;


    public Picture(File file) {
        this.file = file;
    }

    public void initArgs() {
        // 去拓展名
        String name = file.getName().replaceAll("[.][^.]+$", "");
        // 确定输出路径
        String output;
        if (Variables.AUTO_OUTPUT_PATH || Variables.OUTPUT_PATH.isEmpty()) {
            output = String.format("%s/%s/%s%s.%s", file.getParent(),
                    Variables.OUTPUT_PATH_NAME, name, Variables.OUTPUT_PIC_POSTFIX, Variables.OUTPUT_FORMAT);
        } else {
            output = String.format("%s/%s%s.%s", Variables.OUTPUT_PATH, name,
                    Variables.OUTPUT_PIC_POSTFIX, Variables.OUTPUT_FORMAT);
        }
        // 打包参数
        this.args = new String[]{Variables.IMAGE_FLOW_TOOL_PATH, Variables.PROCESS_TYPE, "--in",
                file.getAbsolutePath(), "--out", output, "--command", Variables.command};
    }

    public String[] getArgs() {
        return args;
    }

    public File getFile(){
        return file;
    }
}
