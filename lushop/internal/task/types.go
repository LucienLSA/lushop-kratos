package task

const (
	TASK_WELCOME  = "task:welcome"
	TASK_REMINDER = "task:reminder"
	TASK_PERIODIC = "task:periodic"
)

type TaskPayload struct {
	Username string `yaml:"username" json:"username"`
}
