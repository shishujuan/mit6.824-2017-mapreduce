package mapreduce

import "fmt"
import "sync"

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//

	wp := new(sync.WaitGroup)
	wp.Add(ntasks)

	for i := 0; i < ntasks; i++ {
		go func(taskNumber int) {
			for {
				worker := <-registerChan
				taskArgs := DoTaskArgs{
					JobName:       jobName,
					File:          mapFiles[taskNumber],
					Phase:         phase,
					TaskNumber:    taskNumber,
					NumOtherPhase: n_other,
				}
				ok := call(worker, "Worker.DoTask", taskArgs, nil)
				if ok {
					wp.Done()
					registerChan <- worker
					return
				}
			}
		}(i)
	}

	wp.Wait()
	fmt.Printf("Schedule: %v phase done\n", phase)
}
