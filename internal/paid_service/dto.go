package paid_service

type CreatePaidServiceDTO struct {
	Name         string `json:"name"`
	BaseDuration string `json:"base_duration"`
}
