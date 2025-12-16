package mapreduce

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string,       // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int,             // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Use checkError to handle errors.

	// Clean up any existing merge file for this reduce task to ensure fresh start
	// This prevents contamination from previous test runs
	os.Remove(mergeName(jobName, reduceTaskNumber))
	
	// Map
	kvm := make(map[string][]string)
	// For every nMap (num of file in disk)
	for i := 0; i < nMap; i++{
		// Open each file and read them 
		intermediateFile := reduceName(jobName, i, reduceTaskNumber) // File you will read from
		newFile, err := os.OpenFile(intermediateFile, os.O_RDONLY, 0644)
		if err != nil {
			// log.Fatalf("Failed to open file: %v", err)
			continue
		}
		// Decode the JSON in each file and map the key and multiple values that's goes with that key into a map
		decoder := json.NewDecoder(newFile)
		for {
			var kv KeyValue
			err := decoder.Decode(&kv) // Write stuff to kv
			if err != nil {
				break
			}
			kvm[kv.Key] = append(kvm[kv.Key], kv.Value) // Add the value to the existing key in map
		}
		newFile.Close()
	}
	// Out put file
	mergeFileName := mergeName(jobName, reduceTaskNumber) // Get name of output file
	mergeFile, err := os.OpenFile(mergeFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) // Open output file
	if err != nil {
		log.Fatalf("Failed to open merge file: %v", err)
	}
	enc := json.NewEncoder(mergeFile) // Create encoder
	// For every key in map, call reduceF to get the reduced value for that key
	// Write to output file
	for key, values := range kvm {
		err = enc.Encode(KeyValue{key, reduceF(key, values)}) // Write to output file
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	mergeFile.Close()
}
