public class CompressTask implements Runnable {
    private Picture task;

    public CompressTask(Picture task) {
        this.task = task;
    }

    @Override
    public void run() {
        task.initArgs();
        Main.compress(task);
    }
}