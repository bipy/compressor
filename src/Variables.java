public class Variables {
    // 处理软件IMAGE FLOW位置
    public static final String IMAGE_FLOW_TOOL_PATH =
            "D:\\Users\\Fidelity\\Desktop\\Program\\imageflow.exe";

    // 线程数
    public static final int THREAD_COUNT = 4;

    // 输入路径，可以是文件夹（将递归处理所有子文件夹），如果没有全局utf-8的话要避免中文路径
    public static final String INPUT_PATH =
            "D:\\Users\\Fidelity\\Pictures\\1912X";

    // true: 默认设置
    // false: 在OUTPUT_PATH中设置指定文件夹
    public static Boolean AUTO_OUTPUT_PATH = true;

    // 输出路径，默认为图片当前文件夹下新建compressed文件夹
    public static String OUTPUT_PATH = "";

    // 输出质量 0-100 可选，推荐90
    public static final int QUALITY = 90;

    // 图片编码格式，可以与输出格式不同
    public static final String OUTPUT_FORMAT = "jpg";

    // 是否改变图片大小
    public static final Boolean RESIZE = false;

    // true: 按指定宽度缩放
    // false: 按指定高度缩放
    // 只会缩小，不会放大，RESIZE为false时无效
    public static final Boolean FIXED_WIDTH = true;

    // 指定宽度，RESIZE为false时无效
    public static final int WIDTH = 1920;

    // 指定高度，RESIZE为false时无效
    public static final int HEIGHT = 1080;

    // 输出图片文件名后缀
    public static final String OUTPUT_PIC_POSTFIX = "compressed";

    // 输出图片文件夹名称
    public static final String OUTPUT_PATH_NAME = "compressed";

    // 以下勿修改 ======================================================
    public static final String PROCESS_TYPE = "v0.1/ir4";

    public static String command = "format=" + Variables.OUTPUT_FORMAT + "&quality=" + Variables.QUALITY;

    private Variables(){}
}
