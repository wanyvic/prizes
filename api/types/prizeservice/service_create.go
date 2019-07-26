package prizeservice

import (
	prizestypes "github.com/wanyvic/prizes/api/types"
)

type ServiceCreate struct {
	ServiceName  string `json:"service_name,omitempty"`
	Image        string `json:"image,omitempty"`
	SSHPubkey    string `json:"ssh_pubkey,omitempty"`
	Amount       int64  `json:"amount,omitempty"`
	Pubkey       string `json:"pubkey,omitempty"`
	OutPoint     string `json:"out_point,omitempty"`
	BlockHeight  int64  `json:"block_height,omitempty"`
	ServicePrice int64  `json:"service_price,omitempty"`
	prizestypes.Hardware
	ENV               map[string]string `json:"env,omitempty"`
	MasterNodeN2NAddr string            `json:"master_node_n2n_addr,omitempty"`
}
