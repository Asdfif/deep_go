package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	Tasks map[int]*Task
}

func NewScheduler() Scheduler {
	return Scheduler{Tasks: make(map[int]*Task)}
}

func (s *Scheduler) AddTask(task *Task) {
	s.Tasks[task.Identifier] = task
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if task, ok := s.Tasks[taskID]; ok {
		task.Priority = newPriority
	}
}

func (s *Scheduler) GetTask() Task {
	found := Task{}

	for _, task := range s.Tasks {
		if task.Priority > found.Priority {
			found = *task
		}

		fmt.Printf("FOUND %v\n", found)
	}

	delete(s.Tasks, found.Identifier)
	return found
}

func TestTrace(t *testing.T) {
	task1 := &Task{Identifier: 1, Priority: 10}
	task2 := &Task{Identifier: 2, Priority: 20}
	task3 := &Task{Identifier: 3, Priority: 30}
	task4 := &Task{Identifier: 4, Priority: 40}
	task5 := &Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	fmt.Printf("1 %+v\n", task)
	assert.Equal(t, *task5, task)

	task = scheduler.GetTask()
	fmt.Printf("2 %+v\n", task)

	assert.Equal(t, *task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	fmt.Printf("3 %+v\n", task)

	task = scheduler.GetTask()
	fmt.Printf("4 %+v\n", task)
	assert.Equal(t, *task1, task)

	task = scheduler.GetTask()
	fmt.Printf("5 %+v\n", task)
	assert.Equal(t, *task3, task)
}
