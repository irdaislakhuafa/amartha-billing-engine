package scheduller

import (
	"context"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/usecase"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/go-sdk/log"
)

type (
	Interface interface {
		Run()
		Close() error
	}

	scheduller struct {
		log  log.Interface
		cfg  config.Config
		uc   *usecase.Usecase
		cron gocron.Scheduler
	}
)

var once *sync.Once = &sync.Once{}

func Init(log log.Interface, cfg config.Config, uc *usecase.Usecase) Interface {
	return &scheduller{
		log:  log,
		cfg:  cfg,
		uc:   uc,
		cron: nil,
	}
}

func (s *scheduller) Run() {
	var err error
	once.Do(func() {
		s.cron, err = gocron.NewScheduler()
		if err != nil {
			panic(err)
		}

		s.Register()
		s.cron.Start()
	})
}

func (s *scheduller) Register() {
	// print hello world every 1 second
	s.cron.NewJob(gocron.DurationJob(time.Second), gocron.NewTask(func() {
		s.log.Debug(context.Background(), "Hello World")
	}))
}

func (s *scheduller) Close() error {
	return s.cron.Shutdown()
}
