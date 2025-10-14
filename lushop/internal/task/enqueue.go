package task

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

func NewClient(opt asynq.RedisClientOpt) *asynq.Client {
	return asynq.NewClient(opt)
}

func EnqueueWelcome(c *asynq.Client, username string) (*asynq.TaskInfo, error) {
	b, _ := json.Marshal(TaskPayload{Username: username})
	return c.Enqueue(asynq.NewTask(TASK_WELCOME, b))
}

func EnqueueReminderIn(c *asynq.Client, username string, d time.Duration) (*asynq.TaskInfo, error) {
	b, _ := json.Marshal(TaskPayload{Username: username})
	return c.Enqueue(asynq.NewTask(TASK_REMINDER, b), asynq.ProcessIn(d))
}

// 延迟任务
func EnqueueReminderAdvanced(c *asynq.Client, username string) (*asynq.TaskInfo, error) {
	b, _ := json.Marshal(TaskPayload{Username: username})
	t := asynq.NewTask(TASK_REMINDER, b)
	return c.Enqueue(
		t,
		asynq.ProcessIn(2*time.Hour),  // 延时 2 小时
		asynq.MaxRetry(5),             // 最多重试 5 次
		asynq.Unique(30*time.Minute),  // 30 分钟内去重
		asynq.Queue("default"),        // 指定队列
		asynq.Timeout(30*time.Second), // 单次尝试超时
		asynq.Deadline(time.Now().Add(3*time.Hour)), // 超过截止不再执行
		asynq.Retention(24*time.Hour),               // 结果保留 24 小时
	)
}
