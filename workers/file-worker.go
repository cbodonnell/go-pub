package workers

import (
	"fmt"
	"log"
	"time"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/repositories"
)

type FileWorker struct {
	conf    config.Configuration
	repo    repositories.Repository
	channel chan interface{}
}

func NewFileWorker(_conf config.Configuration, _repo repositories.Repository) Worker {
	return &FederationWorker{
		conf:    _conf,
		repo:    _repo,
		channel: make(chan interface{}),
	}
}

func (f *FileWorker) Start() {
	go func() {
		for {
			err := f.repo.PurgeUnusedFiles()
			if err != nil {
				log.Println(fmt.Sprintf("Failed to purge unused files: %s", err.Error()))
			}
			time.Sleep(24 * time.Hour)
		}
	}()
}

func (f *FileWorker) GetChannel() chan interface{} {
	return f.channel
}
