package service

import (
	prizestypes "github.com/wanyvic/prizes/api/types"
)

type ServiceCreate struct {
	ServiceCreateID      string               `json:"service_create_id,omitempty"`
	ServiceName          string               `json:"service_name,omitempty"`
	Image                string               `json:"image,omitempty"`
	SSHPubkey            string               `json:"ssh_pubkey,omitempty"`
	Amount               int64                `json:"amount,omitempty"`
	Pubkey               string               `json:"pubkey,omitempty"`
	OutPoint             string               `json:"out_point,omitempty"`
	ServicePrice         int64                `json:"service_price,omitempty"`
	Hardware             prizestypes.Hardware `json:"hardware,omitempty"`
	ENV                  map[string]string    `json:"env,omitempty"`
	MasterNodeN2NAddr    string               `json:"master_node_n2n_addr,omitempty"`
	Drawee               string               `json:"drawee,omitempty"`
	MasterNodeFeeRate    int64                `json:"master_node_fee_rate,omitempty"`
	DevFeeRate           int64                `json:"dev_fee_rate,omitempty"` //max 10000
	MasterNodeFeeAddress string               `json:"master_node_fee_address,omitempty"`
	DevFeeAddress        string               `json:"dev_fee_address,omitempty"`
}
