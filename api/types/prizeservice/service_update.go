package prizeservice

type ServiceUpdate struct {
	ServiceID    string `json:"service_id,omitempty"`
	Amount       int64  `json:"amount,omitempty"`
	Pubkey       string `json:"pubkey,omitempty"`
	BlockHeight  int64  `json:"block_height,omitempty"`
	ServicePrice int64  `json:"service_price,omitempty"`
	OutPoint     string `json:"out_point,omitempty"`
}
