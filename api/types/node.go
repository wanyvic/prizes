package types

const (
	LabelRevenueAddress = "REVENUE_ADDRESS"
	LabelCPUType        = "CPUTYPE"
	LabelCPUThread      = "CPUTHREAD"
	LabelMemoryType     = "MEMORYTYPE"
	LabelMemoryCount    = "MEMORYCOUNT"
	LabelGPUType        = "GPUTYPE"
	LabelGPUCount       = "MEMORYCOUNT"
	LabelNFSIP          = "NFSIP"
)

type NodeListStatistics struct {
	WorkerToken       string     `json:"worker_token,omitempty"`
	TotalCount        int        `json:"total_count,omitempty"`
	AvailabilityCount int        `json:"availability_count,omitempty"`
	UsableCount       int        `json:"usable_count,omitempty"`
	List              []NodeInfo `json:"list,omitempty"`
}
type NodeInfo struct {
	NodeID       string            `json:"node_id,omitempty"`
	NodeState    string            `json:"noed_state,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	ReachAddress string            `json:"reach_address,omitempty"`
	Hardware     Hardware          `json:"hardware,omitempty"`
	OnWorking    bool              `json:"onworking,omitempty"`
}
