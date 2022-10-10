// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    rPCActivityForPort8020, err := UnmarshalRPCActivityForPort8020(bytes)
//    bytes, err = rPCActivityForPort8020.Marshal()

package hadoop

import "encoding/json"

func UnmarshalRPCActivityForPort8020(data []byte) (RPCActivityForPort8020, error) {
	var r RPCActivityForPort8020
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RPCActivityForPort8020) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RPCActivityForPort8020 struct {
	Beans []RPCActivityForPort8020Bean `json:"beans,omitempty"`
}

type RPCActivityForPort8020Bean struct {
	Name                             *string  `json:"name,omitempty"`
	ModelerType                      *string  `json:"modelerType,omitempty"`
	TagPort                          *string  `json:"tag.port,omitempty"`
	TagContext                       *string  `json:"tag.Context,omitempty"`
	TagNumOpenConnectionsPerUser     *string  `json:"tag.NumOpenConnectionsPerUser,omitempty"`
	TagHostname                      *string  `json:"tag.Hostname,omitempty"`
	ReceivedBytes                    *int64   `json:"ReceivedBytes,omitempty"`
	SentBytes                        *int64   `json:"SentBytes,omitempty"`
	RPCQueueTimeNumOps               *int64   `json:"RpcQueueTimeNumOps,omitempty"`
	RPCQueueTimeAvgTime              *float64 `json:"RpcQueueTimeAvgTime,omitempty"`
	RPCProcessingTimeNumOps          *int64   `json:"RpcProcessingTimeNumOps,omitempty"`
	RPCProcessingTimeAvgTime         *float64 `json:"RpcProcessingTimeAvgTime,omitempty"`
	DeferredRPCProcessingTimeNumOps  *int64   `json:"DeferredRpcProcessingTimeNumOps,omitempty"`
	DeferredRPCProcessingTimeAvgTime *float32 `json:"DeferredRpcProcessingTimeAvgTime,omitempty"`
	RPCAuthenticationFailures        *int64   `json:"RpcAuthenticationFailures,omitempty"`
	RPCAuthenticationSuccesses       *int64   `json:"RpcAuthenticationSuccesses,omitempty"`
	RPCAuthorizationFailures         *int64   `json:"RpcAuthorizationFailures,omitempty"`
	RPCAuthorizationSuccesses        *int64   `json:"RpcAuthorizationSuccesses,omitempty"`
	RPCClientBackoff                 *int64   `json:"RpcClientBackoff,omitempty"`
	RPCSlowCalls                     *int64   `json:"RpcSlowCalls,omitempty"`
	NumOpenConnections               *int64   `json:"NumOpenConnections,omitempty"`
	CallQueueLength                  *int64   `json:"CallQueueLength,omitempty"`
	NumDroppedConnections            *int64   `json:"NumDroppedConnections,omitempty"`
}
