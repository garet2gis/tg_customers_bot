package service

import "time"

type PaidService struct {
	ID           string        `json:"service_id"`
	Name         string        `json:"name"`
	BaseDuration time.Duration `json:"base_duration"`
}
