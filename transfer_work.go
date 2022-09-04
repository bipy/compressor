package main

// transfer
// when process finished, failCh will be closed
func transferFailList() {
	defer wg.Done()
	for t := range failCh {
		failList = append(failList, t)
	}
}

// transfer
// close the channel cuz this is the only sender
func transferTaskList() {
	defer close(inCh)
	for t := range taskList {
		inCh <- taskList[t]
	}
}
