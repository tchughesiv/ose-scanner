package controller

import (
	"fmt"
)

type Worker struct {
	id         	int
	jobQueue 	chan Job
	workerPool 	chan chan Job
	quit 		chan bool
}

func NewWorker (index int, workerPool chan chan Job) Worker {
	return Worker {
		id:		index,
		workerPool:	workerPool,
		jobQueue: 	make(chan Job),
		quit:		make(chan bool),
	}
}

func (w Worker) Start() {
	fmt.Printf ("Starting worker %d\n", w.id)
	
	go func() {
		for {
			w.workerPool <- w.jobQueue

			select {
				case job := <-w.jobQueue:
					if err := job.ScanImage.scan(); err != nil {
						fmt.Printf("Error scanning image: %s", err.Error())
					}
					job.Done()

				case <-w.quit:
					// we have received a signal to stop
					fmt.Printf ("Aborting worker %d\n", w.id)
					return
			}
		}
	}()
}

