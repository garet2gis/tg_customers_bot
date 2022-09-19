package appointment

type Appointment struct {
	ClientID        string `json:"client_id"`
	ServiceID       string `json:"service_id"`
	AppointmentTime string `json:"appointment_time"`
}
