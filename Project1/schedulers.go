// Intro to OS- Project 1 - Schedulers
// EUID: 11469987
package main

import (
	"fmt"
	"io"
	"sort"
)

type (
	Process struct {
		ProcessID     string
		ArrivalTime   int64
		BurstDuration int64
		Priority      int64
	}
	TimeSlice struct {
		PID   string
		Start int64
		Stop  int64
	}
)

//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
// Print a simple string

// function for implementing FCFS schedule
func FCFSSchedule(w io.Writer, title string, processes []Process) {
	//defining variables
	var (
		//intializing variables
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)

}

// function for implementing SJF schedule
func SJFSchedule(w io.Writer, title string, processes []Process) {
	//defining variables
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	// Sort processes by burst duration (shortest first).
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].BurstDuration < processes[j].BurstDuration
	})

	for i := range processes {
		if processes[i].ArrivalTime > serviceTime {
			// If the process arrives after the current time, update the service time.
			serviceTime = processes[i].ArrivalTime
		}
		// calculating the totalWaitingTime
		waitingTime = serviceTime - processes[i].ArrivalTime
		totalWait += float64(waitingTime)
		// starting
		starting := waitingTime + processes[i].ArrivalTime
		//turnaround
		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)
		//completion
		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)
		// acess ith element of array i.e schedule in this case
		schedule[i] = []string{ // []string to initialize new slice of strings
			fmt.Sprint(processes[i].ProcessID), //fmt.sprint takes different types of data and constructs a string,  without converting each value to string.
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		//adding i-th process's duration to total serviceTime and assigning it back to serviceTime
		serviceTime += processes[i].BurstDuration
		//scheduling tasks
		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID, //task ID
			Start: starting,               //start time of task
			Stop:  serviceTime,            //end time
		})
	}

	totalCount := float64(len(processes)) //counts the elements in processes
	averageThroughput := totalCount / lastCompletion
	averageTurnaround := totalTurnaround / totalCount
	averageWait := totalWait / totalCount

	outputTitle(w, title)
	outputGantt(w, gantt)
	//output schedule
	outputSchedule(w, schedule, averageWait, averageTurnaround, averageThroughput)
}

// function for implementing SJF schedule
func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
	//defining variables for SJF schedule
	var (
		//intializing variables to integer
		serviceTime int64
		waitingTime int64
		//intializing variables to float
		totalTurnaround float64
		totalWait       float64
		lastCompletion  float64
		//initializing schedule with result of // Intializing processes with slice of strings
		schedule = make([][]string, len(processes))
		gantt    = make([]TimeSlice, 0) //creates empty slice, TimeSLice
	)

	// Sort processes by priority (highest priority first), and for processes with the same priority, sort by burst duration (shortest first).
	sort.Slice(processes, func(i, j int) bool {
		if processes[i].Priority == processes[j].Priority {
			return processes[i].BurstDuration < processes[j].BurstDuration
		}
		return processes[i].Priority > processes[j].Priority
	})

	for i := range processes {
		if processes[i].ArrivalTime > serviceTime {
			// If the process arrives after the current time, update service time.
			serviceTime = processes[i].ArrivalTime
		}

		waitingTime = serviceTime - processes[i].ArrivalTime
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration
		// append to gantt with timeslice values (PID, Star, stop) of processes
		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

// function for implementing RRS schedule
func RRSchedule(w io.Writer, title string, processes []Process) {
	//defining variables
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule              = make([][]string, len(processes))
		gantt                 = make([]TimeSlice, 0)
		quantum         int64 = 4 // Adjust the time quantum as needed
	)
	// declare the slice for queue
	queue := make([]Process, len(processes))
	copy(queue, processes)
	// loop
	for len(queue) > 0 {
		process := queue[0]
		queue = queue[1:]

		// If the process arrives after the current time, update service time.
		if process.ArrivalTime > serviceTime {

			serviceTime = process.ArrivalTime
		}
		// Process can't finish in one quantum, so update its properties.
		if process.BurstDuration > quantum {

			waitingTime = serviceTime - process.ArrivalTime
			totalWait += float64(waitingTime)

			schedule = append(schedule, []string{
				fmt.Sprint(process.ProcessID),
				fmt.Sprint(process.Priority),
				fmt.Sprint(quantum),
				fmt.Sprint(process.ArrivalTime),
				fmt.Sprint(waitingTime),
				fmt.Sprint(quantum + waitingTime),
				fmt.Sprint(serviceTime + quantum),
			})

			gantt = append(gantt, TimeSlice{
				PID:   process.ProcessID,
				Start: serviceTime,
				Stop:  serviceTime + quantum,
			})

			serviceTime += quantum

			// Reinsert the process at the end of the queue with reduced burst time.
			process.BurstDuration -= quantum
			queue = append(queue, process)
		} else {
			// Process can finish in this quantum.
			waitingTime = serviceTime - process.ArrivalTime
			totalWait += float64(waitingTime)

			schedule = append(schedule, []string{
				fmt.Sprint(process.ProcessID),
				fmt.Sprint(process.Priority),
				fmt.Sprint(process.BurstDuration),
				fmt.Sprint(process.ArrivalTime),
				fmt.Sprint(waitingTime),
				fmt.Sprint(process.BurstDuration + waitingTime),
				fmt.Sprint(serviceTime + process.BurstDuration),
			})

			gantt = append(gantt, TimeSlice{
				PID:   process.ProcessID,
				Start: serviceTime,
				Stop:  serviceTime + process.BurstDuration,
			})

			serviceTime += process.BurstDuration
		}
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)

}

//endregion
