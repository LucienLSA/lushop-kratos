package task

import (
	"context"

	"lushop/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
)

type AsynqComponent struct {
	logger log.Logger

	srv   *asynq.Server
	sched *asynq.Scheduler
	mgr   *asynq.PeriodicTaskManager
}

func (a *AsynqComponent) Start(ctx context.Context) error {
	if a.srv != nil {
		if err := a.srv.Start(newMux()); err != nil {
			return err
		}
	}
	if a.sched != nil {
		go func() {
			_ = a.sched.Run()
		}()
	}
	if a.mgr != nil {
		go func() {
			_ = a.mgr.Run()
		}()
	}
	return nil
}

func (a *AsynqComponent) Stop(ctx context.Context) error {
	if a.mgr != nil {
		a.mgr.Shutdown()
	}
	if a.sched != nil {
		a.sched.Shutdown()
	}
	if a.srv != nil {
		a.srv.Shutdown()
	}
	return nil
}

type periodicFileConfig struct {
	Configs []struct {
		Cronspec string      `yaml:"cronspec"`
		TaskType string      `yaml:"task_type"`
		Payload  interface{} `yaml:"payload"`
	} `yaml:"configs"`
}

func NewAsynqComponent(data *conf.Data, task *conf.Task, logger log.Logger) (*AsynqComponent, func(), error) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     data.GetRedis().GetAddr(),
		Password: data.GetRedis().GetPassword(),
		DB:       int(data.GetRedis().GetDb()),
	}

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 20,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})
	scheduler := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{})

	// 可选：动态定时任务（从本地 YAML 文件，路径按需调整）
	var mgr *asynq.PeriodicTaskManager

	cleanup := func() {}

	return &AsynqComponent{
		logger: logger,
		srv:    srv,
		sched:  scheduler,
		mgr:    mgr,
	}, cleanup, nil
}
