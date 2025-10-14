package task

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

func newMux() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.Use(loggingMiddleware)

	mux.HandleFunc(TASK_WELCOME, sendWelcomeEmail)
	mux.HandleFunc(TASK_REMINDER, sendReminderEmail)
	mux.HandleFunc(TASK_PERIODIC, sendPeriodicEmail)
	return mux
}

func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		log.Printf("Start processing %q", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			return err
		}
		log.Printf("Finished processing %q", t.Type())
		return nil
	})
}

func sendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var p TaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	// TODO: 实际发送逻辑
	log.Printf("[WELCOME] username=%s", p.Username)
	return nil
}

func sendReminderEmail(ctx context.Context, t *asynq.Task) error {
	var p TaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	// TODO: 实际发送逻辑
	log.Printf("[REMINDER] username=%s", p.Username)
	return nil
}

func sendPeriodicEmail(ctx context.Context, t *asynq.Task) error {
	var p TaskPayload
	if len(t.Payload()) > 0 {
		_ = json.Unmarshal(t.Payload(), &p)
	}
	// TODO: 实际发送逻辑
	log.Printf("[PERIODIC] username=%s", p.Username)
	return nil
}
