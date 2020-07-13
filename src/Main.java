import java.io.*;
import java.util.ArrayList;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {
    // 运行时
    private static Runtime runtime = Runtime.getRuntime();

    // 计数
    private static int sum = 1;
    private static AtomicInteger count = new AtomicInteger(0);

    // 处理图片集合
    private static ArrayList<Picture> picList = new ArrayList<>();
    private static ArrayList<Picture> failList = new ArrayList<>();


    public static void main(String[] args) {
        init();
        File inputFile = new File(Variables.INPUT_PATH);
        if (inputFile.isDirectory()) {
            find(inputFile);
            sum = picList.size();
            process();
        } else if (inputFile.isFile()) {
            process(new Picture(inputFile));
        }
        // 统计
        if (!failList.isEmpty()) {
            System.out.println("Oops! Some of them are failed:");
            for (Picture f : failList) {
                System.out.println("Fail: " + f.getFile().getPath());
            }
        }
        System.out.println(String.format("\nProcess Complete! Total: %d - Success: %d - Fail: %d",
                sum, sum - failList.size(), failList.size()));
    }

    public static void init() {
        // 检查参数是否合法
        try {
            if (!new File(Variables.IMAGE_FLOW_TOOL_PATH).exists() && !new File(Variables.IMAGE_FLOW_TOOL_PATH).canExecute()) {
                throw new Exception("ERROR: imageflow_tool NOT FOUND");
            }
            if (!Variables.AUTO_OUTPUT_PATH && new File(Variables.OUTPUT_PATH).isFile()) {
                throw new Exception("ERROR: OUTPUT PATH SHOULD BE A DIRECTORY");
            }
            if (!new File(Variables.INPUT_PATH).exists()) {
                throw new Exception("ERROR: INPUT FILE NOT FOUND");
            }
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }
        if (Variables.RESIZE) {
            if (Variables.FIXED_WIDTH) {
                Variables.command += "&width=" + Variables.WIDTH;
            } else {
                Variables.command += "&height=" + Variables.HEIGHT;
            }
        }
        if (Variables.OVERWRITE) {
            Variables.OUTPUT_PATH_NAME = "";
            Variables.OUTPUT_PIC_POSTFIX = ".compress_temp";
        }
    }

    public static void find(File currentFile) {
        // 递归访问文件夹，并将所有图片放入集合
        File[] files = currentFile.listFiles();
        ArrayList<Picture> currentPicList = new ArrayList<>();
        for (File file : files) {
            // 防止递归处理已压缩的图片
            if (file.isDirectory() && !file.getName().equals(Variables.OUTPUT_PATH_NAME)) {
                find(file);
            } else if (file.getName().toLowerCase().matches(".*[.](png|jpg|jpge)$")) {
                currentPicList.add(new Picture(file));
            }
        }
        // 判断输出文件夹是否存在
        if (!currentPicList.isEmpty()) {
            File outputPath;
            if (Variables.AUTO_OUTPUT_PATH || Variables.OUTPUT_PATH.isEmpty()) {
                outputPath = new File(currentFile.getPath() + "/" + Variables.OUTPUT_PATH_NAME);
            } else {
                outputPath = new File(Variables.OUTPUT_PATH);
            }
            if (!outputPath.exists()) {
                outputPath.mkdir();
            }
        }
        picList.addAll(currentPicList);
    }

    public static void process(Picture inputFile) {
        System.out.println("======= Single File Mode =======");
        if (Variables.AUTO_OUTPUT_PATH || Variables.OUTPUT_PATH.isEmpty()) {
            Variables.OUTPUT_PATH = inputFile.getFile().getParent();
            Variables.AUTO_OUTPUT_PATH = false;
        }
        inputFile.initArgs();
        compress(inputFile);
    }

    public static void process() {
        // 多线程处理
        if (Variables.THREAD_COUNT > 1) {
            System.out.println("======= Multi Thread Mode =======");
            ExecutorService service = Executors.newFixedThreadPool(Variables.THREAD_COUNT);
            for (Picture pic : picList) {
                service.submit(new CompressTask(pic));
            }
            service.shutdown();
            try {
                service.awaitTermination(72, TimeUnit.HOURS);
            } catch (InterruptedException e) {
                System.out.println("运行超时");
                System.exit(2);
            }
        } else {
            System.out.println("======= Single Thread Mode =======");
            for (Picture pic : picList) {
                pic.initArgs();
                compress(pic);
            }
        }
    }

    public static Boolean overwrite(File source, File temp) {
        return source.delete() && temp.renameTo(
                new File(source.getPath().replaceAll("[^.]+$", Variables.OUTPUT_FORMAT)));
    }

    public static void compress(Picture pic) {
        if (compress(pic.getArgs()) && (!Variables.OVERWRITE || overwrite(pic.getFile(), new File(pic.getArgs()[5])))) {
            System.out.println(String.format("(%d/%d) %s success",
                    count.addAndGet(1), sum, pic.getFile().getPath()));
            return;
        }
        failList.add(pic);
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
                    br.close();
                    return true;
                }
            }
            br.close();
            return false;
        } catch (Exception e) {
            return false;
        }
    }
}

