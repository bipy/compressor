import java.io.*;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.concurrent.*;

public class Main {
    // 处理软件IMAGE FLOW位置
    private static final String IMAGE_FLOW_TOOL_PATH =
            "D:\\Users\\Fidelity\\Desktop\\Program\\imageflow.exe";

    // 线程数
    private static final int THREAD_COUNT = 4;

    // 线程池
    private static ExecutorService service;

    // 处理图片集合
    private static ArrayList<File> picList = new ArrayList<>();

    // 输入路径，可以是文件夹（将递归处理所有子文件夹），如果没有全局utf-8的话要避免中文路径
    private static final String INPUT_PATH =
            "D:\\Users\\Fidelity\\Pictures";

    // 输出路径，默认为图片当前文件夹下新建compressed文件夹
    private static String OUTPUT_PATH = "";

    // true: 默认设置
    // false: 在OUTPUT_PATH中设置指定文件夹
    private static Boolean AUTO_OUTPUT_PATH = true;

    // 输出质量 0-100 可选，推荐90
    private static final int QUALITY = 90;

    // 图片编码格式，可以与输出格式不同
    private static final String OUTPUT_FORMAT = "jpg";

    // 是否改变图片大小
    private static final Boolean RESIZE = false;

    // true: 按指定宽度缩放
    // false: 按指定高度缩放
    // 只会缩小，不会放大，RESIZE为false时无效
    private static final Boolean FIXED_WIDTH = true;

    // 指定宽度，RESIZE为false时无效
    private static final int WIDTH = 1920;

    // 指定高度，RESIZE为false时无效
    private static final int HEIGHT = 1080;

    // 输出图片文件名后缀
    private static final String OUTPUT_PIC_POSTFIX = "compressed";

    // 输出图片文件夹名称
    private static final String OUTPUT_PATH_NAME = "compressed";

    // 以下勿修改
    private static final String PROCESS_TYPE = "v0.1/ir4";
    private static String command = "format=" + OUTPUT_FORMAT + "&quality=" + QUALITY;
    private static Runtime runtime = Runtime.getRuntime();

    public static void main(String[] args) {
        if(!new File(IMAGE_FLOW_TOOL_PATH).exists()){
            System.out.println("ERROR: Imageflow_tool NOT FOUND");
            return;
        }
        if (RESIZE) setReSize();
        File inputFile = new File(INPUT_PATH);
        if (inputFile.isDirectory()) {
            if (!AUTO_OUTPUT_PATH && new File(OUTPUT_PATH).isFile()) {
                System.out.println("ERROR: Output Path should be a Directory.");
                return;
            }
            find(inputFile);
            // 多线程处理
            if (THREAD_COUNT > 1) {
                System.out.println("======= Multi Thread Mode =======");
                service = Executors.newFixedThreadPool(THREAD_COUNT);
                File[] picFiles = picList.toArray(new File[picList.size()]);
                int batch = picFiles.length / THREAD_COUNT;
                for (int i = 0; i < THREAD_COUNT - 1; i++) {
                    int startPos = i * batch;
                    int endPos = (i + 1) * batch;
                    service.submit(new CompressTask(Arrays.copyOfRange(picFiles, startPos, endPos)));
                }
                // 最后一组处理剩下的所有图片
                service.submit(new CompressTask(Arrays.copyOfRange(picFiles, (THREAD_COUNT - 1) * batch, picFiles.length)));
                service.shutdown();
            } else {
                System.out.println("======= Single Thread Mode =======");
                for (File pic : picList) {
                    initArgs(pic);
                }
            }
        } else if (inputFile.isFile()) {
            System.out.println("======= Single File Mode =======");
            if (AUTO_OUTPUT_PATH || OUTPUT_PATH.isEmpty()) {
                OUTPUT_PATH = inputFile.getParent();
                AUTO_OUTPUT_PATH = false;
            }
            initArgs(inputFile);
        }

    }

    public static void find(File currentFile) {
        // 递归访问文件夹，并将所有图片放入集合
        File[] files = currentFile.listFiles();
        ArrayList<File> currentPicList = new ArrayList<>();
        for (File file : files) {
            if (file.isDirectory()) {
                find(file);
            } else if (file.getName().matches(".*[.](png|jpg|jpge)$")) {
                currentPicList.add(file);
            }
        }
        // 判断输出文件夹是否存在
        if (!currentPicList.isEmpty()) {
            File outputPath;
            if (AUTO_OUTPUT_PATH || OUTPUT_PATH.isEmpty()) {
                outputPath = new File(currentFile.getPath() + "\\" + OUTPUT_PATH_NAME);
            } else {
                outputPath = new File(OUTPUT_PATH);
            }
            if (!outputPath.exists()) {
                outputPath.mkdir();
            }
        }
        picList.addAll(currentPicList);
    }

    public static void setReSize() {
        if (FIXED_WIDTH) {
            command += "&width=" + WIDTH;
        } else {
            command += "&height=" + HEIGHT;
        }
    }

    public static void initArgs(File pic) {
        // 去拓展名
        String name = pic.getName().replaceAll("[.][^.]+$", "");
        // 确定输出路径
        String output;
        if (AUTO_OUTPUT_PATH || OUTPUT_PATH.isEmpty()) {
            output = String.format("%s\\%s\\%s_%s.%s", pic.getParent(),
                    OUTPUT_PATH_NAME, name, OUTPUT_PIC_POSTFIX, OUTPUT_FORMAT);
        } else {
            output = String.format("%s\\%s_%s.%s", OUTPUT_PATH, name,
                    OUTPUT_PIC_POSTFIX, OUTPUT_FORMAT);
        }
        // 打包参数
        String[] args = new String[]{IMAGE_FLOW_TOOL_PATH, PROCESS_TYPE, "--in",
                pic.getAbsolutePath(), "--out", output, "--command", command};
        if (!compress(args)) {
            System.out.println(pic.getAbsolutePath() + " fail");
        } else {
            System.out.println(pic.getAbsolutePath() + " success");
        }
    }


    public static Boolean compress(String[] args) {
        try {
            Process p = runtime.exec(args);
            BufferedReader br = new BufferedReader(new InputStreamReader(p.getInputStream()));
            String line;
            while ((line = br.readLine()) != null) {
                // 判断是否成功
                if (line.contains("200")) {
                    if (p.isAlive()) {
                        p.destroy();
                    }
                    return true;
                }
            }
            return false;
        } catch (IOException e) {
            return false;
        }
    }
}



