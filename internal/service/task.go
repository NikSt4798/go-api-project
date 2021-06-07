package service

import (
	"context"
	"etokubernetes/internal/dto"
	"fmt"
	"sync"
	"time"
)

type TaskService interface {
	GetTask(ctx context.Context) (*dto.Task, error)
	CreateTask(ctx context.Context, task *dto.Task) (*dto.Task, error)
	ReportTask(ctx context.Context, task *dto.Task) error
	List(ctx context.Context) ([]*dto.Task, error)
}

type taskService struct {
	// TODO use repository instead tasks slice
	tasks           []*dto.Task
	taskChangeMutex sync.Mutex
}

func NewTaskService() TaskService {
	return &taskService{tasks: make([]*dto.Task, 0, 1000)}
}

func (s *taskService) GetTask(ctx context.Context) (*dto.Task, error) {
	s.taskChangeMutex.Lock()
	defer s.taskChangeMutex.Unlock()

	var task *dto.Task
	var taskIndex int
	for i, v := range s.tasks {
		if v.Status == dto.TaskStatusNotScheduled || v.Status == dto.TaskStatusTimeout {
			taskIndex = i
			task = v
			task.Status = dto.TaskStatusRunning
			break
		}
	}

	if task != nil {
		task.Timeout = time.Now().Add(60 * time.Second)
		go s.timeoutTask(taskIndex)
		return task, nil
	}
	return nil, fmt.Errorf("task to schedule not found")
}

func (s *taskService) CreateTask(ctx context.Context, task *dto.Task) (*dto.Task, error) {
	// TODO add task validation
	task.ID = len(s.tasks) + 1
	task.Status = dto.TaskStatusNotScheduled

	s.tasks = append(s.tasks, task)

	return task, nil
}

func (s *taskService) ReportTask(ctx context.Context, task *dto.Task) error {
	s.taskChangeMutex.Lock()
	defer s.taskChangeMutex.Unlock()

	if len(s.tasks) < task.ID {
		return fmt.Errorf("task not found")
	}

	taskIndex := task.ID - 1

	s.tasks[taskIndex].Status = task.Status
	s.tasks[taskIndex].Result = task.Result

	return nil
}

func (s *taskService) List(ctx context.Context) ([]*dto.Task, error) {
	return s.tasks, nil
}

func (s *taskService) timeoutTask(taskIndex int) {
	if len(s.tasks) < taskIndex {
		return
	}

	task := s.tasks[taskIndex]

	time.Sleep(task.Timeout.Sub(time.Now()))

	s.taskChangeMutex.Lock()
	defer s.taskChangeMutex.Unlock()

	task = s.tasks[taskIndex]

	if task.Status == dto.TaskStatusRunning {
		s.tasks[taskIndex].Status = dto.TaskStatusTimeout
	}
}
