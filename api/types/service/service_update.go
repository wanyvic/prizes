package service

type ServiceUpdate struct {
	ServiceUpdateID string `json:"service_update_id,omitempty"`
	ServiceID       string `json:"service_id,omitempty"`
	Amount          int64  `json:"amount,omitempty"`
	Pubkey          string `json:"pubkey,omitempty"`
	ServicePrice    int64  `json:"service_price,omitempty"`
	OutPoint        string `json:"out_point,omitempty"`
	Drawee          string `json:"drawee,omitempty"`
}
