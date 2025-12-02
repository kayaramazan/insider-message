package job

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kayaramazan/insider-message/api/service"
)

type Job struct {
	Mu             sync.RWMutex
	Running        bool
	interval       time.Duration
	cancel         context.CancelFunc
	lastRun        time.Time
	runCount       int64
	messageService *service.MessageService
}

func New(interval time.Duration, messageService *service.MessageService) *Job {
	return &Job{
		interval:       interval,
		messageService: messageService,
	}
}

func (job *Job) Start() {
	job.Mu.Lock()
	if job.Running {
		job.Mu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	job.cancel = cancel
	job.Running = true
	job.Mu.Unlock()

	go job.run(ctx)
	log.Println("Job started")
}

func (j *Job) Stop() {
	j.Mu.Lock()
	defer j.Mu.Unlock()

	if !j.Running {
		return
	}

	j.cancel()
	j.Running = false
	log.Println("Job stopped")
}

func (j *Job) run(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	j.doWork()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			j.doWork()
		}
	}
}

func (j *Job) doWork() {
	j.Mu.Lock()
	j.lastRun = time.Now()
	j.runCount++
	j.Mu.Unlock()
	err := j.messageService.SendMessage(context.Background())
	if err != nil {
		log.Println("Messages could not sending ", err)
	}
}

func (j *Job) IsRunning() bool {
	j.Mu.RLock()
	defer j.Mu.RUnlock()
	return j.Running
}

func (j *Job) Toggle() {
	j.Mu.Lock()
	running := j.Running
	j.Mu.Unlock()

	if running {
		j.Stop()
	} else {
		j.Start()
	}
}
