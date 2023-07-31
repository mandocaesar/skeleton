package event

import "encoding/json"

type SampleEvent struct {
	ID string `json:"id"`
}

func (e SampleEvent) Marshal() ([]byte, error) {
	payload, err := json.Marshal(e)

	return payload, err
}
