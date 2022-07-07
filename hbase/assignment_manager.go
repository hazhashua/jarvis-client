// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    assignmentManager, err := UnmarshalAssignmentManager(bytes)
//    bytes, err = assignmentManager.Marshal()

package hbase

import "encoding/json"

func UnmarshalAssignmentManager(data []byte) (AssignmentManager, error) {
	var r AssignmentManager
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AssignmentManager) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AssignmentManager struct {
	Beans []AssignmentManagerBean `json:"beans,omitempty"`
}

type AssignmentManagerBean struct {
	Name                        *string `json:"name,omitempty"`
	ModelerType                 *string `json:"modelerType,omitempty"`
	TagContext                  *string `json:"tag.Context,omitempty"`
	TagHostname                 *string `json:"tag.Hostname,omitempty"`
	OperationCount              *int64  `json:"operationCount,omitempty"`
	DeadServerOpenRegions       *int64  `json:"deadServerOpenRegions,omitempty"`
	RitOldestAge                *int64  `json:"ritOldestAge,omitempty"`
	UnknownServerOpenRegions    *int64  `json:"unknownServerOpenRegions,omitempty"`
	RitDurationNumOps           *int64  `json:"RitDuration_num_ops,omitempty"`
	RitDurationMin              *int64  `json:"RitDuration_min,omitempty"`
	RitDurationMax              *int64  `json:"RitDuration_max,omitempty"`
	RitDurationMean             *int64  `json:"RitDuration_mean,omitempty"`
	RitDuration25ThPercentile   *int64  `json:"RitDuration_25th_percentile,omitempty"`
	RitDurationMedian           *int64  `json:"RitDuration_median,omitempty"`
	RitDuration75ThPercentile   *int64  `json:"RitDuration_75th_percentile,omitempty"`
	RitDuration90ThPercentile   *int64  `json:"RitDuration_90th_percentile,omitempty"`
	RitDuration95ThPercentile   *int64  `json:"RitDuration_95th_percentile,omitempty"`
	RitDuration98ThPercentile   *int64  `json:"RitDuration_98th_percentile,omitempty"`
	RitDuration99ThPercentile   *int64  `json:"RitDuration_99th_percentile,omitempty"`
	RitDuration999ThPercentile  *int64  `json:"RitDuration_99.9th_percentile,omitempty"`
	RitCount                    *int64  `json:"ritCount,omitempty"`
	RitCountOverThreshold       *int64  `json:"ritCountOverThreshold,omitempty"`
	MoveSubmittedCount          *int64  `json:"MoveSubmittedCount,omitempty"`
	SplitFailedCount            *int64  `json:"SplitFailedCount,omitempty"`
	CloseSubmittedCount         *int64  `json:"CloseSubmittedCount,omitempty"`
	ReopenSubmittedCount        *int64  `json:"ReopenSubmittedCount,omitempty"`
	OpenFailedCount             *int64  `json:"OpenFailedCount,omitempty"`
	CloseTimeNumOps             *int64  `json:"CloseTime_num_ops,omitempty"`
	CloseTimeMin                *int64  `json:"CloseTime_min,omitempty"`
	CloseTimeMax                *int64  `json:"CloseTime_max,omitempty"`
	CloseTimeMean               *int64  `json:"CloseTime_mean,omitempty"`
	CloseTime25ThPercentile     *int64  `json:"CloseTime_25th_percentile,omitempty"`
	CloseTimeMedian             *int64  `json:"CloseTime_median,omitempty"`
	CloseTime75ThPercentile     *int64  `json:"CloseTime_75th_percentile,omitempty"`
	CloseTime90ThPercentile     *int64  `json:"CloseTime_90th_percentile,omitempty"`
	CloseTime95ThPercentile     *int64  `json:"CloseTime_95th_percentile,omitempty"`
	CloseTime98ThPercentile     *int64  `json:"CloseTime_98th_percentile,omitempty"`
	CloseTime99ThPercentile     *int64  `json:"CloseTime_99th_percentile,omitempty"`
	CloseTime999ThPercentile    *int64  `json:"CloseTime_99.9th_percentile,omitempty"`
	MergeSubmittedCount         *int64  `json:"MergeSubmittedCount,omitempty"`
	CloseFailedCount            *int64  `json:"CloseFailedCount,omitempty"`
	MoveFailedCount             *int64  `json:"MoveFailedCount,omitempty"`
	OpenTimeNumOps              *int64  `json:"OpenTime_num_ops,omitempty"`
	OpenTimeMin                 *int64  `json:"OpenTime_min,omitempty"`
	OpenTimeMax                 *int64  `json:"OpenTime_max,omitempty"`
	OpenTimeMean                *int64  `json:"OpenTime_mean,omitempty"`
	OpenTime25ThPercentile      *int64  `json:"OpenTime_25th_percentile,omitempty"`
	OpenTimeMedian              *int64  `json:"OpenTime_median,omitempty"`
	OpenTime75ThPercentile      *int64  `json:"OpenTime_75th_percentile,omitempty"`
	OpenTime90ThPercentile      *int64  `json:"OpenTime_90th_percentile,omitempty"`
	OpenTime95ThPercentile      *int64  `json:"OpenTime_95th_percentile,omitempty"`
	OpenTime98ThPercentile      *int64  `json:"OpenTime_98th_percentile,omitempty"`
	OpenTime99ThPercentile      *int64  `json:"OpenTime_99th_percentile,omitempty"`
	OpenTime999ThPercentile     *int64  `json:"OpenTime_99.9th_percentile,omitempty"`
	SplitTimeNumOps             *int64  `json:"SplitTime_num_ops,omitempty"`
	SplitTimeMin                *int64  `json:"SplitTime_min,omitempty"`
	SplitTimeMax                *int64  `json:"SplitTime_max,omitempty"`
	SplitTimeMean               *int64  `json:"SplitTime_mean,omitempty"`
	SplitTime25ThPercentile     *int64  `json:"SplitTime_25th_percentile,omitempty"`
	SplitTimeMedian             *int64  `json:"SplitTime_median,omitempty"`
	SplitTime75ThPercentile     *int64  `json:"SplitTime_75th_percentile,omitempty"`
	SplitTime90ThPercentile     *int64  `json:"SplitTime_90th_percentile,omitempty"`
	SplitTime95ThPercentile     *int64  `json:"SplitTime_95th_percentile,omitempty"`
	SplitTime98ThPercentile     *int64  `json:"SplitTime_98th_percentile,omitempty"`
	SplitTime99ThPercentile     *int64  `json:"SplitTime_99th_percentile,omitempty"`
	SplitTime999ThPercentile    *int64  `json:"SplitTime_99.9th_percentile,omitempty"`
	AssignSubmittedCount        *int64  `json:"AssignSubmittedCount,omitempty"`
	MergeTimeNumOps             *int64  `json:"MergeTime_num_ops,omitempty"`
	MergeTimeMin                *int64  `json:"MergeTime_min,omitempty"`
	MergeTimeMax                *int64  `json:"MergeTime_max,omitempty"`
	MergeTimeMean               *int64  `json:"MergeTime_mean,omitempty"`
	MergeTime25ThPercentile     *int64  `json:"MergeTime_25th_percentile,omitempty"`
	MergeTimeMedian             *int64  `json:"MergeTime_median,omitempty"`
	MergeTime75ThPercentile     *int64  `json:"MergeTime_75th_percentile,omitempty"`
	MergeTime90ThPercentile     *int64  `json:"MergeTime_90th_percentile,omitempty"`
	MergeTime95ThPercentile     *int64  `json:"MergeTime_95th_percentile,omitempty"`
	MergeTime98ThPercentile     *int64  `json:"MergeTime_98th_percentile,omitempty"`
	MergeTime99ThPercentile     *int64  `json:"MergeTime_99th_percentile,omitempty"`
	MergeTime999ThPercentile    *int64  `json:"MergeTime_99.9th_percentile,omitempty"`
	UnassignTimeNumOps          *int64  `json:"UnassignTime_num_ops,omitempty"`
	UnassignTimeMin             *int64  `json:"UnassignTime_min,omitempty"`
	UnassignTimeMax             *int64  `json:"UnassignTime_max,omitempty"`
	UnassignTimeMean            *int64  `json:"UnassignTime_mean,omitempty"`
	UnassignTime25ThPercentile  *int64  `json:"UnassignTime_25th_percentile,omitempty"`
	UnassignTimeMedian          *int64  `json:"UnassignTime_median,omitempty"`
	UnassignTime75ThPercentile  *int64  `json:"UnassignTime_75th_percentile,omitempty"`
	UnassignTime90ThPercentile  *int64  `json:"UnassignTime_90th_percentile,omitempty"`
	UnassignTime95ThPercentile  *int64  `json:"UnassignTime_95th_percentile,omitempty"`
	UnassignTime98ThPercentile  *int64  `json:"UnassignTime_98th_percentile,omitempty"`
	UnassignTime99ThPercentile  *int64  `json:"UnassignTime_99th_percentile,omitempty"`
	UnassignTime999ThPercentile *int64  `json:"UnassignTime_99.9th_percentile,omitempty"`
	OpenSubmittedCount          *int64  `json:"OpenSubmittedCount,omitempty"`
	AssignFailedCount           *int64  `json:"AssignFailedCount,omitempty"`
	UnassignSubmittedCount      *int64  `json:"UnassignSubmittedCount,omitempty"`
	MergeFailedCount            *int64  `json:"MergeFailedCount,omitempty"`
	SplitSubmittedCount         *int64  `json:"SplitSubmittedCount,omitempty"`
	ReopenTimeNumOps            *int64  `json:"ReopenTime_num_ops,omitempty"`
	ReopenTimeMin               *int64  `json:"ReopenTime_min,omitempty"`
	ReopenTimeMax               *int64  `json:"ReopenTime_max,omitempty"`
	ReopenTimeMean              *int64  `json:"ReopenTime_mean,omitempty"`
	ReopenTime25ThPercentile    *int64  `json:"ReopenTime_25th_percentile,omitempty"`
	ReopenTimeMedian            *int64  `json:"ReopenTime_median,omitempty"`
	ReopenTime75ThPercentile    *int64  `json:"ReopenTime_75th_percentile,omitempty"`
	ReopenTime90ThPercentile    *int64  `json:"ReopenTime_90th_percentile,omitempty"`
	ReopenTime95ThPercentile    *int64  `json:"ReopenTime_95th_percentile,omitempty"`
	ReopenTime98ThPercentile    *int64  `json:"ReopenTime_98th_percentile,omitempty"`
	ReopenTime99ThPercentile    *int64  `json:"ReopenTime_99th_percentile,omitempty"`
	ReopenTime999ThPercentile   *int64  `json:"ReopenTime_99.9th_percentile,omitempty"`
	MoveTimeNumOps              *int64  `json:"MoveTime_num_ops,omitempty"`
	MoveTimeMin                 *int64  `json:"MoveTime_min,omitempty"`
	MoveTimeMax                 *int64  `json:"MoveTime_max,omitempty"`
	MoveTimeMean                *int64  `json:"MoveTime_mean,omitempty"`
	MoveTime25ThPercentile      *int64  `json:"MoveTime_25th_percentile,omitempty"`
	MoveTimeMedian              *int64  `json:"MoveTime_median,omitempty"`
	MoveTime75ThPercentile      *int64  `json:"MoveTime_75th_percentile,omitempty"`
	MoveTime90ThPercentile      *int64  `json:"MoveTime_90th_percentile,omitempty"`
	MoveTime95ThPercentile      *int64  `json:"MoveTime_95th_percentile,omitempty"`
	MoveTime98ThPercentile      *int64  `json:"MoveTime_98th_percentile,omitempty"`
	MoveTime99ThPercentile      *int64  `json:"MoveTime_99th_percentile,omitempty"`
	MoveTime999ThPercentile     *int64  `json:"MoveTime_99.9th_percentile,omitempty"`
	UnassignFailedCount         *int64  `json:"UnassignFailedCount,omitempty"`
	ReopenFailedCount           *int64  `json:"ReopenFailedCount,omitempty"`
	AssignTimeNumOps            *int64  `json:"AssignTime_num_ops,omitempty"`
	AssignTimeMin               *int64  `json:"AssignTime_min,omitempty"`
	AssignTimeMax               *int64  `json:"AssignTime_max,omitempty"`
	AssignTimeMean              *int64  `json:"AssignTime_mean,omitempty"`
	AssignTime25ThPercentile    *int64  `json:"AssignTime_25th_percentile,omitempty"`
	AssignTimeMedian            *int64  `json:"AssignTime_median,omitempty"`
	AssignTime75ThPercentile    *int64  `json:"AssignTime_75th_percentile,omitempty"`
	AssignTime90ThPercentile    *int64  `json:"AssignTime_90th_percentile,omitempty"`
	AssignTime95ThPercentile    *int64  `json:"AssignTime_95th_percentile,omitempty"`
	AssignTime98ThPercentile    *int64  `json:"AssignTime_98th_percentile,omitempty"`
	AssignTime99ThPercentile    *int64  `json:"AssignTime_99th_percentile,omitempty"`
	AssignTime999ThPercentile   *int64  `json:"AssignTime_99.9th_percentile,omitempty"`
}
