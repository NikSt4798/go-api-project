package dto

import "time"

var TaskStatuses = []string{
	"TaskStatusNotSchedule",
	"TaskStatusRunning",
	"TaskStatusSuccess",
	"TaskStatusError",
	"TaskStatusTimeout",
   }
   
const (
	TaskStatusNotScheduled = int8(iota)
	TaskStatusRunning
	TaskStatusSuccess
	TaskStatusError
	TaskStatusTimeout
)

type Task struct {
	ID         int
	Type       string
	Parameters map[string]interface{}
	Result     interface{}
	Status     int8
	Timeout    time.Time
}
