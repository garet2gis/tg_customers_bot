package paid_service

import (
	"fmt"
	"strings"
	"time"
)

type PaidService struct {
	ID           string        `json:"service_id"`
	Name         string        `json:"name"`
	BaseDuration time.Duration `json:"base_duration"`
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

func toRussian(dur string) string {
	replacer := strings.NewReplacer("d", "д ", "h", "ч ", "m", "м ", "s", "с ")
	return replacer.Replace(dur)
}

func (ps *PaidService) String() string {
	return fmt.Sprintf("Название: %s\nДлительность: %v\n", ps.Name, toRussian(shortDur(ps.BaseDuration)))
}
