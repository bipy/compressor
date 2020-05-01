public class CompressTask implements Runnable {
    private Picture[] taskList;

    public CompressTask(Picture[] taskList) {
        this.taskList = taskList;
    }

    @Override
    public void run() {
        for (Picture pic : taskList) {
            pic.initArgs();
            Main.compress(pic);
        }
    }
}