// Code generated by https://github.com/zhuharev/gen DO NOT EDIT.

package job

import (
	"context"
	"log"
	"time"
)

// Job describe background job
type Job struct {
	Name     string
	Func     func(context.Context) error
	Interval time.Duration
}

// Run runs background jobs
func Run(jobs ...Job) {
	for _, job := range jobs {
		go func(job Job) {
			for ticker := time.NewTicker(job.Interval); ; <-ticker.C {
				err := job.Func(context.Background())
				if err != nil {
					log.Printf("err background job name=%s err=%s", job.Name, err)
				}
			}
		}(job)
	}
}
