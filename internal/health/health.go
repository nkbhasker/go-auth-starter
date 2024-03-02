package health

import (
	"encoding/json"
)

type HealthStatusEnum string

type Checker interface {
	Check() *Health
}

type Health struct {
	status HealthStatusEnum
	info   map[string]interface{}
}

const (
	HealthStatusUp        HealthStatusEnum = "UP"
	HealthStatusDown      HealthStatusEnum = "DOWN"
	HealthStatusUpUnknown HealthStatusEnum = "UNKNOWN"
)

func NewHealth() *Health {
	return &Health{
		status: HealthStatusUpUnknown,
		info:   make(map[string]interface{}),
	}
}

func (h *Health) Status() HealthStatusEnum {
	return h.status
}

func (h *Health) Info() map[string]interface{} {
	data := map[string]interface{}{}
	data["status"] = h.status
	for k, v := range h.info {
		data[k] = v
	}

	return data
}

func (h *Health) SetStatus(status HealthStatusEnum) *Health {
	h.status = status

	return h
}

func (h *Health) SetInfo(key string, value interface{}) *Health {
	h.info[key] = value

	return h
}

func (h *Health) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Info())
}
