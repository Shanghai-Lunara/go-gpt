package operator

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type Worker interface {
	Run()
	Add(t *Task)
}

type worker struct {
	p         Project
	ch        <-chan struct{}
	workQueue workqueue.RateLimitingInterface
}

func (w *worker) Run() {
	defer w.workQueue.ShutDown()
	go wait.Until(w.runWorker, time.Second, w.ch)
	klog.Info("Started workers")
	<-w.ch
	klog.Info("Shutting down workers")
}

func (w *worker) Add(t *Task) {
	w.workQueue.Add(t)
}

func (w *worker) runWorker() {
	for w.processNextWorkItem() {
	}
}

func (w *worker) processNextWorkItem() bool {
	obj, shutdown := w.workQueue.Get()
	if shutdown {
		return false
	}
	_ = func(obj interface{}) error {
		defer w.workQueue.Done(obj)
		var task *Task
		var ok bool
		if task, ok = obj.(*Task); !ok {
			w.workQueue.Forget(obj)
			klog.V(2).Infof("transfer failed t:%v", task)
			return nil
		}
		if err := w.syncHandler(task); err != nil {
			klog.V(2).Info("syncHandler err:", err)
			defer task.ChangeStatus(TaskError)
			//w.workQueue.AddRateLimited(task)
			w.workQueue.Forget(obj)
			return fmt.Errorf("error syncing:%v err:%s, requeuing", task, err.Error())
		}
		w.workQueue.Forget(obj)
		klog.Infof("Successfully synced '%v'", task)
		return nil
	}(obj)
	return true
}

func (w *worker) syncHandler(t *Task) error {
	t.ChangeStatus(TaskProcessing)
	defer t.ChangeStatus(TaskCompleted)
	c := t.Command
	switch c.Command {
	case TaskCmdGitGen:
		return w.p.GitGenerate(c.ProjectName, c.BranchName)
	case TaskCmdSvnCommit:
		return w.p.SvnCommit(c.ProjectName, c.BranchName, c.Message)
	case TaskCmdFtpUpload:
		return w.p.FtpCompress(c.ProjectName, c.BranchName, c.ZipType, c.ZipFlags)
	}
	return nil
}

func NewWorker(ch <-chan struct{}, p Project) Worker {
	var w Worker = &worker{
		p:         p,
		ch:        ch,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "workQueue"),
	}
	go w.Run()
	return w
}
