package paid_service

type CreatePaidServiceDTO struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	BaseDuration string `json:"base_duration"`
}
