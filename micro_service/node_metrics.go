// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    nodeMetrics, err := UnmarshalNodeMetrics(bytes)
//    bytes, err = nodeMetrics.Marshal()

package micro_service

import "encoding/json"

func UnmarshalNodeMetrics(data []byte) (NodeMetrics, error) {
	var r NodeMetrics
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NodeMetrics) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NodeMetrics struct {
	Kind       *string              `json:"kind,omitempty"`
	APIVersion *string              `json:"apiVersion,omitempty"`
	Metadata   *NodeMetricsMetadata `json:"metadata,omitempty"`
	Items      []NodeMetricsItem    `json:"items,omitempty"`
}

type NodeMetricsItem struct {
	Metadata  *NodeMetricsItemMetadata `json:"metadata,omitempty"`
	Timestamp *string                  `json:"timestamp,omitempty"`
	Window    *string                  `json:"window,omitempty"`
	Usage     *Usage                   `json:"usage,omitempty"`
}

type NodeMetricsItemMetadata struct {
	Name              *string `json:"name,omitempty"`
	SelfLink          *string `json:"selfLink,omitempty"`
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`
}

type Usage struct {
	CPU    *string `json:"cpu,omitempty"`
	Memory *string `json:"memory,omitempty"`
}

type NodeMetricsMetadata struct {
	SelfLink *string `json:"selfLink,omitempty"`
}
