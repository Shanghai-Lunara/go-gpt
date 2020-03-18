package operator

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TaskWaiting = iota
	TaskProcessing
	TaskCompleted
	TaskError
)

const (
	TaskCmdGitGen    = "gitGen"
	TaskCmdSvnCommit = "svnCommit"
	TaskCmdFtpUpload = "ftpUpload"
)

type Task struct {
	mu sync.RWMutex

	Id      int      `json:"id"`
	Status  int      `json:"status"`
	Message []string `json:"message"`
	Command Command  `json:"command"`
}

func (t *Task) ChangeStatus(status int) {
	t.mu.Lock()
	t.Status = status
	var msg string
	var cot = map[int]string{
		TaskWaiting:    "waiting",
		TaskProcessing: "processing",
		TaskCompleted:  "completed",
		TaskError:      "error",
	}
	t.mu.Unlock()
	msg = fmt.Sprintf("[%s] change status: %s", time.Now().Format("2006-01-02 15:04:05"), cot[status])
	t.AppendMessage(msg)
}

func (t *Task) AppendMessage(msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Message = append(t.Message, msg)
}

type TaskHub struct {
	mu    sync.RWMutex
	Max   int32
	Tasks map[int]*Task
}

func (th *TaskHub) getNextTaskId() int {
	return int(atomic.AddInt32(&th.Max, 1))
}

func (th *TaskHub) NewTask(c *Command) *Task {
	t := &Task{
		Id:      th.getNextTaskId(),
		Status:  TaskWaiting,
		Message: make([]string, 0),
		Command: *c,
	}
	th.mu.Lock()
	defer th.mu.Unlock()
	th.Tasks[t.Id] = t
	return t
}

func (th *TaskHub) GetAll() map[int]Task {
	th.mu.RLock()
	defer th.mu.RUnlock()
	res := make(map[int]Task, len(th.Tasks))
	for k, v := range th.Tasks {
		v.mu.RLock()
		res[k] = *v
		v.mu.RUnlock()
	}
	return res
}

func NewTaskHub() *TaskHub {
	return &TaskHub{
		Max:   0,
		Tasks: make(map[int]*Task, 0),
	}
}
