package job

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kayaramazan/insider-message/api/service"
)

type Job interface {
	IsRunning() bool
	Toggle()
	Start()
}

type jobImpl struct {
	mu             sync.RWMutex
	Running        bool
	interval       time.Duration
	cancel         context.CancelFunc
	lastRun        time.Time
	runCount       int64
	messageService service.MessageService
}

func New(interval time.Duration, messageService service.MessageService) Job {
	return &jobImpl{
		interval:       interval,
		messageService: messageService,
	}
}

func (job *jobImpl) Start() {
	job.mu.Lock()
	if job.Running {
		job.mu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	job.cancel = cancel
	job.Running = true
	job.mu.Unlock()

	go job.run(ctx)
	log.Println("Job started")
}

func (j *jobImpl) stop() {
	j.mu.Lock()
	defer j.mu.Unlock()

	if !j.Running {
		return
	}

	j.cancel()
	j.Running = false
	log.Println("Job stopped")
}

func (j *jobImpl) run(ctx context.Context) {
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

func (j *jobImpl) doWork() {
	j.mu.Lock()
	j.lastRun = time.Now()
	j.runCount++
	j.mu.Unlock()
	err := j.messageService.SendMessage(context.Background())
	if err != nil {
		log.Println("Messages could not sending ", err)
	}
}

func (j *jobImpl) IsRunning() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.Running
}

func (j *jobImpl) Toggle() {
	j.mu.Lock()
	running := j.Running
	j.mu.Unlock()

	if running {
		j.stop()
	} else {
		j.Start()
	}
}
