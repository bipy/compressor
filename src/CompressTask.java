import java.io.File;

public class CompressTask implements Runnable {
    private File[] taskList;

    public CompressTask(File[] taskList) {
        this.taskList = taskList;
    }

    @Override
    public void run() {
        for (File pic : taskList) {
            Main.initArgs(pic);
        }
    }
}