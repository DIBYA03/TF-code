package workpool

import (
	"log"
	"sync"
)

//Task is a task for a giving worker
type Task struct {
	//The error return by the task
	//This error is store so can access if needed
	Err error

	//The function to run concurrently
	f func() error
}

func (t *Task) start(wg *sync.WaitGroup) {
	t.Err = t.f()
	wg.Done()
}

//NewTask return a new task
func NewTask(f func() error) *Task {
	return &Task{f: f}
}

//Pool the worker pool
type Pool interface {
	Run()
	Errors() []error
	AddTasks([]*Task)
	AddOne(*Task)
}

type pool struct {
	Tasks       []*Task
	Concurrency int
	TaskChan    chan *Task
	wg          sync.WaitGroup
}

//The work task for a single routine
func (p *pool) work() {
	for task := range p.TaskChan {
		task.start(&p.wg)
	}
}

//New return a new worker pool
func New(tasks []*Task, concurrency int) Pool {
	return &pool{
		Tasks:       tasks,
		Concurrency: concurrency,
		TaskChan:    make(chan *Task),
	}
}

//NewEmpty will return an empty Pool
//Youll need to set the `tasks`
func NewEmpty() Pool {
	return &pool{TaskChan: make(chan *Task), Concurrency: 10}
}

//Run runs all the tasks
func (p *pool) Run() {
	if len(p.Tasks) == 0 {
		log.Print("No task to run at this time")
		return
	}

	log.Printf("Running %v task(s) at %v.", len(p.Tasks), p.Concurrency)
	for i := 0; i < p.Concurrency; i++ {
		go p.work()
	}
	//Add tasks to wait group
	p.wg.Add(len(p.Tasks))

	for _, task := range p.Tasks {
		p.TaskChan <- task
	}
	//All workers are done
	close(p.TaskChan)

	p.wg.Wait()
}

func (p *pool) AddTasks(ts []*Task) {
	for _, task := range ts {
		p.Tasks = append(p.Tasks, task)
	}
}

func (p *pool) AddOne(t *Task) {
	p.Tasks = append(p.Tasks, t)
}

//Errors All the errors return by tasks
func (p *pool) Errors() []error {
	var errs []error
	for _, task := range p.Tasks {
		if task.Err != nil {
			errs = append(errs, task.Err)
		}
	}
	return errs
}
