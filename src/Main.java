import java.io.*;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {
    // 运行时
    private static Runtime runtime = Runtime.getRuntime();

    // 计数
    private static int sum = 1;
    private static AtomicInteger count = new AtomicInteger(0);

    // 处理图片集合
    private static List<Picture> picList = new ArrayList<>();
    private static List<String> failList = Collections.synchronizedList(new ArrayList<>());


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
            for (String f : failList) {
                System.out.println("Fail: " + f);
            }
        }
        System.out.println(String.format("\nProcess Complete! Total: %d - Success: %d - Fail: %d",
                sum, sum - failList.size(), failList.size()));
    }

    private static void init() {
        // 检查参数是否合法
        try {
            if (!new File(Variables.IMAGE_FLOW_TOOL_PATH).exists() || !new File(Variables.IMAGE_FLOW_TOOL_PATH).canExecute()) {
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

    private static void find(File currentFile) {
        // 递归访问文件夹，并将所有图片放入集合
        File[] files = currentFile.listFiles();
        ArrayList<Picture> currentPicList = new ArrayList<>();
        for (File file : files) {
            // 防止递归处理已压缩的图片
            if (file.isDirectory() && !file.getName().equals(Variables.OUTPUT_PATH_NAME)) {
                find(file);
            } else if (file.getName().toLowerCase().matches(".*[.](png|jpg)$")) {
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

    private static void process(Picture inputFile) {
        System.out.println("======= Single File Mode =======");
        if (Variables.AUTO_OUTPUT_PATH || Variables.OUTPUT_PATH.isEmpty()) {
            Variables.OUTPUT_PATH = inputFile.getFile().getParent();
            Variables.AUTO_OUTPUT_PATH = false;
        }
        compress(inputFile);
    }

    private static void process() {
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
                compress(pic);
            }
        }
    }

    private static void appendFailList(Picture pic) {
        failList.add(pic.getInputPath());
    }

    private static Boolean overwrite(File source, File temp) {
        return source.delete() && temp.renameTo(
                new File(source.getPath().replaceAll("[^.]+$", Variables.OUTPUT_FORMAT)));
    }

    static void compress(Picture pic) {
        // 打包参数
        String[] args = new String[]{
                Variables.IMAGE_FLOW_TOOL_PATH,
                Variables.PROCESS_TYPE,
                "--in",
                pic.getInputPath(),
                "--out",
                pic.getOutputPath(),
                "--command",
                Variables.command
        };

        if(!pic.getFile().exists()){
            appendFailList(pic);
            return;
        }
        if (compress(args)) {
            if (Variables.OVERWRITE) {
                if (!overwrite(pic.getFile(), new File(pic.getOutputPath()))) {
                    appendFailList(pic);
                    return;
                }
            }
            System.out.println(String.format("(%d/%d) %s succeed",
                    count.incrementAndGet(), sum, pic.getFile().getPath()));
            return;
        }
        appendFailList(pic);
    }

    private static Boolean compress(String[] args) {
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

class CompressTask implements Runnable {
    private Picture task;

    CompressTask(Picture task) {
        this.task = task;
    }

    @Override
    public void run() {
        Main.compress(task);
    }
}

class Picture {
    private File file;
    private String outputPath;
    private String inputPath;


    Picture(File file) {
        this.file = file;
        this.inputPath = file.getPath();
        // 去拓展名
        String name = file.getName().replaceAll("[.][^.]+$", "");
        // 确定输出路径
        if (Variables.AUTO_OUTPUT_PATH || Variables.OUTPUT_PATH.isEmpty()) {
            outputPath = String.format("%s/%s/%s%s.%s", file.getParent(),
                    Variables.OUTPUT_PATH_NAME, name, Variables.OUTPUT_PIC_POSTFIX, Variables.OUTPUT_FORMAT);
        } else {
            outputPath = String.format("%s/%s%s.%s", Variables.OUTPUT_PATH, name,
                    Variables.OUTPUT_PIC_POSTFIX, Variables.OUTPUT_FORMAT);
        }
    }

    String getOutputPath() {
        return outputPath;
    }

    String getInputPath() {
        return inputPath;
    }

    File getFile() {
        return file;
    }
}


