package request

// Simplified webhook payload
type RundeckWebhook struct {
	ExecutionID int       `json:"executionId"`
	Status      string    `json:"status"`
	Execution   execution `json:"execution"`
}

type execution struct {
	Job job `json:"job"`
}

type job struct {
	ID      string                 `json:"id"`
	Options map[string]interface{} `json:"options"`
}
