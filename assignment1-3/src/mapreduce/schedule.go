package mapreduce

import "log"

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

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//

	taskIndex := 0
	workers := mr.workers
	go func() {
		for {
			worker, isOpen := <-mr.registerChannel
			debug("worker %s, isOpen %t\n", worker, isOpen)
			workers = append(workers, worker)
		}
	}()

	for {
		if taskIndex == ntasks {
			break
		}
		for _, worker := range workers {
			debug("here\n")
			arg := DoTaskArgs{
				JobName:       mr.jobName,
				File:          mr.files[taskIndex],
				Phase:         phase,
				TaskNumber:    taskIndex,
				NumOtherPhase: nios,
			}
			ok := call(worker, "Worker.DoTask", arg, new(struct{}))
			if !ok {
				log.Fatal("Failed to schedule")
			}
			taskIndex++
			if taskIndex == ntasks {
				break
			}
		}

		debug("Task number %d", taskIndex)
	}

	debug("Schedule: %v phase done\n", phase)
}
