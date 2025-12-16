# MapReduce Project in Go

This project implements a **MapReduce framework** in Go, supporting both **sequential** and **distributed execution**. It includes handling worker failures and demonstrates fault-tolerant distributed computation.

---

## Overview

- Build a **MapReduce library** in Go.  
- Support **sequential execution** for debugging and **parallel execution** across multiple workers.  
- Implement **map** and **reduce** functions for example applications (`dgrep` and `dgrepcount`).  
- Handle **worker failures** gracefully in the distributed setup.  

---

## Features

### MapReduce Flow

1. **Input**: List of files, map function, reduce function, number of reduce tasks.  
2. **Master**: Creates tasks, starts RPC server, assigns tasks to workers, handles failures.  
3. **doMap()**: Reads input file, applies map function, writes intermediate key/value pairs to `nReduce` files.  
4. **doReduce()**: Collects intermediate files for each reduce task, applies reduce function, writes output.  
5. **Merge**: Combines all reduce outputs into a final result file.  
6. **Shutdown**: Sends RPC to terminate workers.

---

### Applications

1. **dgrep.go**: Prints lines in input files matching a specified pattern.  
   - Example output matches:
     ```bash
     grep -H -n <pattern> <files>
     ```
   - Each key/value pair corresponds to a matching line in a file.

2. **dgrepcount.go**: Counts the number of lines containing a specified pattern.  
   - Example output matches:
     ```bash
     grep -c <pattern> <files>
     ```
   - Each key/value pair corresponds to the number of matching lines per file.

---

### Execution Modes

- **Sequential**: Executes all map tasks, then all reduce tasks, one at a time. Useful for debugging.  
  ```bash
  go run dgrep.go master sequential <pattern> <files>
  ```
- **Distributed**: Executes map and reduce tasks concurrently across worker threads using RPC.
  ```bash
  go run dgrep.go master distributed <pattern> <files>
  ```

### Worker Fault Tolerance
- If a worker fails during a task, the master reassigns the task to another available worker.
- Multiple executions of the same task are safe due to the functional nature of map and reduce functions.

### Implementation Notes
- Modify mapF() and reduceF() in dgrep.go and dgrepcount.go for application logic.
- Modify schedule() in mapreduce/schedule.go to assign tasks to workers efficiently.
- Use Go channels, goroutines, and sync.WaitGroup for concurrency.
- Use RPC for communication between master and workers.
- Intermediate files follow the naming convention:
- mrtmp.<job>-res-<map/reduce task number>
- Output files:
  - mrtmp.dgrep for dgrep.go
  - mrtmp.dgrepcount for dgrepcount.go

### Testing
- Sequential tests:
  ```bash
  cd src/mapreduce
  go test -run Sequential
  ```
- Distributed tests:
  ```bash
  go test -run TestBasic
  ```
- Failure handling tests:
  ```bash
  go test -run Failure
  ```
- Full automated test:
  ```bash
  sh ./test-mr.sh
  ```

### Learning Outcomes
- Understand distributed computation and fault tolerance.
- Apply concurrency concepts in Go (goroutines, channels, WaitGroups).
- Implement MapReduce applications for real-world text processing.
- Learn to debug distributed systems and handle failures gracefully.
