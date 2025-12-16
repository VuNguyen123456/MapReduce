package mapreduce

import(
	"sync"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// (mr *Master) is the receiver that make schedule a method of any Master object
	// Basically it can access Master field with mr.registerChannel

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	
	// ok := call(master, "Master.Register", args, new(struct{}))
	var wg sync.WaitGroup
    for i := 0; i < ntasks; i++ {
		// debug("schedule: launching task %d\n", i)
		wg.Add(1)
		//Get a avai worker from registerChan
		//Send task to worker with RPC parallel
		go func(taskNum int) {
			defer wg.Done() // Mark this task as done when the goroutine completes
			for {
				workerAddr := <-mr.registerChannel //get a worker from channel
				// Prepare DoTaskArgs
				args := DoTaskArgs{
					JobName: mr.jobName,
					InputArg: mr.arg,
					Phase: phase,
					TaskNumber: taskNum,
					NumOtherPhase: nios,
				}
				if phase == mapPhase{
					args.File = mr.files[taskNum]
				}
				

				var reply struct{}
				ok := call(workerAddr, "Worker.DoTask", &args, &reply)
				

				if ok{
					//if success, put the worker back to the channel
					// mr.registerChannel <- workerAddr
					go func(w string) { mr.registerChannel <- w }(workerAddr) //Sent asynchronously
					break
				}else{
					//if fail, reschedule the task by sending same task to a different worker
					// Assign imediately to another available worker
					continue
				}
			}
		}(i)
	}
	wg.Wait() // Wait for all tasks to complete
}
