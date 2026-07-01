package sshreversetunnel

import "github.com/gravitational/teleport/api/types"

// Taken from Teleport codebase
type DialReq struct {
	// Address is the target host to make a connection to.
	Address string `json:"address,omitempty"`

	// ServerID is the hostUUID.clusterName of the node. ServerID is used when
	// dialing through a tunnel to SSH and application nodes.
	ServerID string `json:"server_id,omitempty"`

	// ConnType is the type of connection requested, either node or application.
	ConnType types.TunnelType `json:"conn_type"`

	// ClientSrcAddr is the original observed client address, it is used to propagate
	// correct client IP through indirect connections inside teleport
	ClientSrcAddr string `json:"client_src_addr,omitempty"`

	// ClientDstAddr is the original client's destination address, it is used to propagate
	// correct client point of contact through indirect connections inside teleport
	ClientDstAddr string `json:"client_dst_addr,omitempty"`

	// IsAgentlessNode specifies whether the target is an agentless node.
	IsAgentlessNode bool `json:"is_agentless_node,omitempty"`
}
