package service

type ServiceUpdate struct {
	ServiceUpdateID      string `json:"service_update_id,omitempty"`
	ServiceID            string `json:"service_id,omitempty"`
	Amount               int64  `json:"amount"`
	Pubkey               string `json:"pubkey,omitempty"`
	ServicePrice         int64  `json:"service_price"`
	OutPoint             string `json:"out_point,omitempty"`
	Drawee               string `json:"drawee,omitempty"`
	MasterNodeFeeRate    int64  `json:"master_node_fee_rate"`
	DevFeeRate           int64  `json:"dev_fee_rate"` //max 10000
	MasterNodeFeeAddress string `json:"master_node_fee_address,omitempty"`
	DevFeeAddress        string `json:"dev_fee_address,omitempty"`
}
