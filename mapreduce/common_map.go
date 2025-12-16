package mapreduce

import (
	"hash/fnv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	input_argument string, //the input argument for the map function
	inFile string,    // name of the input file
	nReduce int, // the number of reduce tasks that will be run ("R" in the paper)
	mapF func(file string, arg string, contents string) []KeyValue,
) {
	// TODO:
	// You will need to write this function.
	// You can find the filename for this map task's input to reduce task number
	// r using reduceName(jobName, mapTaskNumber, r). The ihash function (given
	// below doMap) should be used to decide which file a given key belongs into.
	//
	// The intermediate output of a map task is stored in the file
	// system as multiple files whose name indicates which map task produced
	// them, as well as which reduce task they are for. Coming up with a
	// scheme for how to store the key/value pairs on disk can be tricky,
	// especially when taking into account that both keys and values could
	// contain newlines, quotes, and any other character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!
	// Use checkError to handle errors.
	
	// Clean up any existing intermediate files for this map task to ensure fresh start
	// This prevents contamination from previous test runs
	for r := 0; r < nReduce; r++ {
		os.Remove(reduceName(jobName, mapTaskNumber, r))
	}
	
	// Read the input file
	data, err := os.ReadFile(inFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	contents := string(data) // string

	// Call map Function so it can process that work and produce intermediate result
	processes := mapF(inFile, input_argument, contents) // This will produce a list of KeyValue: both are string

	// ihash the keys so we know which file given key belong too to get it to reduce
	// Reduce name to get the name of the file
	// Create nReduce files with correct name
	// Write key/value pairs into the correct file in JSON format
	for i := 0; i < len(processes); i++ {
		r := int(ihash(processes[i].Key)) % nReduce
		name := reduceName(jobName, mapTaskNumber, r)

		newFile, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}

		jsonData := json.NewEncoder(newFile)
		err = jsonData.Encode(&processes[i])
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		newFile.Close()
	}
}



func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
