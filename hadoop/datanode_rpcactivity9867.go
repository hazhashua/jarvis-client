// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    rPCActivityForPort9867, err := UnmarshalRPCActivityForPort9867(bytes)
//    bytes, err = rPCActivityForPort9867.Marshal()

package hadoop

import "encoding/json"

func UnmarshalRPCActivityForPort9867(data []byte) (RPCActivityForPort9867, error) {
	var r RPCActivityForPort9867
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RPCActivityForPort9867) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RPCActivityForPort9867 struct {
	Beans []RpcActivityForPort9867Bean `json:"beans,omitempty"`
}

type RpcActivityForPort9867Bean struct {
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
	RPCProcessingTimeAvgTime         *int64   `json:"RpcProcessingTimeAvgTime,omitempty"`
	DeferredRPCProcessingTimeNumOps  *int64   `json:"DeferredRpcProcessingTimeNumOps,omitempty"`
	DeferredRPCProcessingTimeAvgTime *int64   `json:"DeferredRpcProcessingTimeAvgTime,omitempty"`
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
