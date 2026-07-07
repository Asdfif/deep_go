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
	Heap      []*Task
	Positions map[int]int
}

func NewScheduler() Scheduler {
	return Scheduler{Heap: []*Task{}, Positions: make(map[int]int)}
}

func (s *Scheduler) ShiftUpFrom(index int) {
	for taskIdentifier, id := range s.Positions {
		if id >= index {
			s.Positions[taskIdentifier] += 1
		}
	}
}

func (s *Scheduler) ShiftDownFrom(index int) {
	for taskIdentifier, id := range s.Positions {
		if id >= index {
			s.Positions[taskIdentifier] -= 1
		}
	}
}

func (s *Scheduler) AddTask(task *Task) {
	fmt.Printf("ADD %v\n", task)
	if len(s.Heap) == 0 {
		s.Heap = append(s.Heap, task)
		s.Positions[task.Identifier] = 0
		return
	}

	if s.Heap[len(s.Heap)-1].Priority <= task.Priority {
		s.Heap = append(s.Heap, task)
		s.Positions[task.Identifier] = len(s.Heap)
		return
	}

	if s.Heap[0].Priority >= task.Priority {
		s.Heap = append([]*Task{task}, s.Heap...)
		s.ShiftUpFrom(0)
		s.Positions[task.Identifier] = 0
		return
	}

	i := len(s.Heap) / 2

	for {
		fmt.Printf("TASK %v; I = %v\n", task, i)
		if (s.Heap[i-1].Priority < task.Priority && s.Heap[i].Priority > task.Priority) || s.Heap[i-1].Priority == task.Priority {
			fmt.Printf("TASK %+v s.Heap[%v].Priority=%v;\n", task, i, s.Heap[i].Priority)

			s.Heap = append(s.Heap[:i], append([]*Task{task}, s.Heap[i:]...)...)
			s.ShiftUpFrom(i)
			s.Positions[task.Identifier] = i
			return
		}

		if s.Heap[i].Priority > task.Priority {
			i = i / 2
			continue
		} else {
			i = i + (len(s.Heap)-i)/2
			continue
		}
	}
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	position := s.Positions[taskID]
	task := s.Heap[position]

	if position == 0 {
		s.Heap = s.Heap[1:]
	} else if position == len(s.Heap)-1 {
		s.Heap = s.Heap[:position]
	} else {
		s.Heap = append(s.Heap[:position], s.Heap[position+1:]...)
	}

	task.Priority = newPriority
	s.AddTask(task)
}

func (s *Scheduler) Print() {
	fmt.Print("PRINT SCHEDULER ")
	for i, task := range s.Heap {
		fmt.Printf("%v) %v; ", i, *task)
	}
	fmt.Println()
}

func (s *Scheduler) GetTask() Task {
	found := *s.Heap[len(s.Heap)-1]

	s.Heap = s.Heap[:len(s.Heap)-1]

	return found
}

func TestTrace(t *testing.T) {
	task10 := &Task{Identifier: 1, Priority: 10}
	task20 := &Task{Identifier: 2, Priority: 20}
	task30 := &Task{Identifier: 3, Priority: 30}
	task40 := &Task{Identifier: 4, Priority: 40}
	task50 := &Task{Identifier: 5, Priority: 50}
	task60 := &Task{Identifier: 6, Priority: 60}
	task70 := &Task{Identifier: 7, Priority: 70}
	task62 := &Task{Identifier: 8, Priority: 62}
	task61 := &Task{Identifier: 9, Priority: 61}

	scheduler := NewScheduler()
	scheduler.AddTask(task10)
	scheduler.AddTask(task20)
	scheduler.AddTask(task30)
	scheduler.AddTask(task40)
	scheduler.AddTask(task50)

	scheduler.AddTask(task60)
	scheduler.AddTask(task70)
	scheduler.AddTask(task62)
	scheduler.AddTask(task61)
	scheduler.Print()

	task := scheduler.GetTask() // &Task{Identifier: 7, Priority: 70}
	assert.Equal(t, *task70, task)

	task = scheduler.GetTask() // &Task{Identifier: 8, Priority: 62}
	assert.Equal(t, *task62, task)

	scheduler.ChangeTaskPriority(9, 53) // &Task{Identifier: 9, Priority: 61} -> &Task{Identifier: 9, Priority: 53}
	scheduler.Print()

	task = scheduler.GetTask() // &Task{Identifier: 6, Priority: 60}
	assert.Equal(t, *task60, task)

	task = scheduler.GetTask() // &Task{Identifier: 9, Priority: 53}
	assert.Equal(t, *task61, task)
}
