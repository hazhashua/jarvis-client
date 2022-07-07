// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    clusterMetrics, err := UnmarshalClusterMetrics(bytes)
//    bytes, err = clusterMetrics.Marshal()

package hadoop

import "encoding/json"

func UnmarshalClusterMetrics(data []byte) (ClusterMetrics, error) {
	var r ClusterMetrics
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ClusterMetrics) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ClusterMetrics struct {
	Beans []ClusterMetricBean `json:"beans,omitempty"`
}

type ClusterMetricBean struct {
	Name                   *string `json:"name,omitempty"`
	ModelerType            *string `json:"modelerType,omitempty"`
	TagClusterMetrics      *string `json:"tag.ClusterMetrics,omitempty"`
	TagContext             *string `json:"tag.Context,omitempty"`
	TagHostname            *string `json:"tag.Hostname,omitempty"`
	NumActiveNMS           *int64  `json:"NumActiveNMs,omitempty"`
	NumDecommissioningNMS  *int64  `json:"NumDecommissioningNMs,omitempty"`
	NumDecommissionedNMS   *int64  `json:"NumDecommissionedNMs,omitempty"`
	NumLostNMS             *int64  `json:"NumLostNMs,omitempty"`
	NumUnhealthyNMS        *int64  `json:"NumUnhealthyNMs,omitempty"`
	NumRebootedNMS         *int64  `json:"NumRebootedNMs,omitempty"`
	NumShutdownNMS         *int64  `json:"NumShutdownNMs,omitempty"`
	AMLaunchDelayNumOps    *int64  `json:"AMLaunchDelayNumOps,omitempty"`
	AMLaunchDelayAvgTime   *int64  `json:"AMLaunchDelayAvgTime,omitempty"`
	AMRegisterDelayNumOps  *int64  `json:"AMRegisterDelayNumOps,omitempty"`
	AMRegisterDelayAvgTime *int64  `json:"AMRegisterDelayAvgTime,omitempty"`
}
