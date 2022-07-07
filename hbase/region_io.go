// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    regionserverIO, err := UnmarshalRegionserverIO(bytes)
//    bytes, err = regionserverIO.Marshal()

package hbase

import "encoding/json"

func UnmarshalRegionserverIO(data []byte) (RegionserverIO, error) {
	var r RegionserverIO
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RegionserverIO) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RegionserverIO struct {
	Beans []IOBean `json:"beans,omitempty"`
}

type IOBean struct {
	Name                       *string `json:"name,omitempty"`
	ModelerType                *string `json:"modelerType,omitempty"`
	TagContext                 *string `json:"tag.Context,omitempty"`
	TagHostname                *string `json:"tag.Hostname,omitempty"`
	FSChecksumFailureCount     *int64  `json:"fsChecksumFailureCount,omitempty"`
	FSPReadTimeNumOps          *int64  `json:"FsPReadTime_num_ops,omitempty"`
	FSPReadTimeMin             *int64  `json:"FsPReadTime_min,omitempty"`
	FSPReadTimeMax             *int64  `json:"FsPReadTime_max,omitempty"`
	FSPReadTimeMean            *int64  `json:"FsPReadTime_mean,omitempty"`
	FSPReadTime25ThPercentile  *int64  `json:"FsPReadTime_25th_percentile,omitempty"`
	FSPReadTimeMedian          *int64  `json:"FsPReadTime_median,omitempty"`
	FSPReadTime75ThPercentile  *int64  `json:"FsPReadTime_75th_percentile,omitempty"`
	FSPReadTime90ThPercentile  *int64  `json:"FsPReadTime_90th_percentile,omitempty"`
	FSPReadTime95ThPercentile  *int64  `json:"FsPReadTime_95th_percentile,omitempty"`
	FSPReadTime98ThPercentile  *int64  `json:"FsPReadTime_98th_percentile,omitempty"`
	FSPReadTime99ThPercentile  *int64  `json:"FsPReadTime_99th_percentile,omitempty"`
	FSPReadTime999ThPercentile *int64  `json:"FsPReadTime_99.9th_percentile,omitempty"`
	FSWriteTimeNumOps          *int64  `json:"FsWriteTime_num_ops,omitempty"`
	FSWriteTimeMin             *int64  `json:"FsWriteTime_min,omitempty"`
	FSWriteTimeMax             *int64  `json:"FsWriteTime_max,omitempty"`
	FSWriteTimeMean            *int64  `json:"FsWriteTime_mean,omitempty"`
	FSWriteTime25ThPercentile  *int64  `json:"FsWriteTime_25th_percentile,omitempty"`
	FSWriteTimeMedian          *int64  `json:"FsWriteTime_median,omitempty"`
	FSWriteTime75ThPercentile  *int64  `json:"FsWriteTime_75th_percentile,omitempty"`
	FSWriteTime90ThPercentile  *int64  `json:"FsWriteTime_90th_percentile,omitempty"`
	FSWriteTime95ThPercentile  *int64  `json:"FsWriteTime_95th_percentile,omitempty"`
	FSWriteTime98ThPercentile  *int64  `json:"FsWriteTime_98th_percentile,omitempty"`
	FSWriteTime99ThPercentile  *int64  `json:"FsWriteTime_99th_percentile,omitempty"`
	FSWriteTime999ThPercentile *int64  `json:"FsWriteTime_99.9th_percentile,omitempty"`
	FSReadTimeNumOps           *int64  `json:"FsReadTime_num_ops,omitempty"`
	FSReadTimeMin              *int64  `json:"FsReadTime_min,omitempty"`
	FSReadTimeMax              *int64  `json:"FsReadTime_max,omitempty"`
	FSReadTimeMean             *int64  `json:"FsReadTime_mean,omitempty"`
	FSReadTime25ThPercentile   *int64  `json:"FsReadTime_25th_percentile,omitempty"`
	FSReadTimeMedian           *int64  `json:"FsReadTime_median,omitempty"`
	FSReadTime75ThPercentile   *int64  `json:"FsReadTime_75th_percentile,omitempty"`
	FSReadTime90ThPercentile   *int64  `json:"FsReadTime_90th_percentile,omitempty"`
	FSReadTime95ThPercentile   *int64  `json:"FsReadTime_95th_percentile,omitempty"`
	FSReadTime98ThPercentile   *int64  `json:"FsReadTime_98th_percentile,omitempty"`
	FSReadTime99ThPercentile   *int64  `json:"FsReadTime_99th_percentile,omitempty"`
	FSReadTime999ThPercentile  *int64  `json:"FsReadTime_99.9th_percentile,omitempty"`
}
