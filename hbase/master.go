package hbase

import "encoding/json"

func UnmarshalMasterMain(data []byte) (MasterMain, error) {
	var r MasterMain
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MasterMain) MarshalMasterMain() ([]byte, error) {
	return json.Marshal(r)
}

type MasterMain struct {
	Beans []MasterBean `json:"beans,omitempty"`
}

type MasterBean struct {
	Name                             *string  `json:"name,omitempty"`
	ModelerType                      *string  `json:"modelerType,omitempty"`
	TagLiveRegionServers             *string  `json:"tag.liveRegionServers,omitempty"`
	TagDeadRegionServers             *string  `json:"tag.deadRegionServers,omitempty"`
	TagDraininigRegionServers        *string  `json:"tag.draininigRegionServers,omitempty"`
	TagZookeeperQuorum               *string  `json:"tag.zookeeperQuorum,omitempty"`
	TagServerName                    *string  `json:"tag.serverName,omitempty"`
	TagClusterID                     *string  `json:"tag.clusterId,omitempty"`
	TagIsActiveMaster                *string  `json:"tag.isActiveMaster,omitempty"`
	TagContext                       *string  `json:"tag.Context,omitempty"`
	TagHostname                      *string  `json:"tag.Hostname,omitempty"`
	MergePlanCount                   *int64   `json:"mergePlanCount,omitempty"`
	SplitPlanCount                   *int64   `json:"splitPlanCount,omitempty"`
	MasterActiveTime                 *int64   `json:"masterActiveTime,omitempty"`
	MasterStartTime                  *int64   `json:"masterStartTime,omitempty"`
	MasterFinishedInitializationTime *int64   `json:"masterFinishedInitializationTime,omitempty"`
	AverageLoad                      *float32 `json:"averageLoad,omitempty"`
	NumRegionServers                 *int64   `json:"numRegionServers,omitempty"`
	NumDeadRegionServers             *int64   `json:"numDeadRegionServers,omitempty"`
	NumDrainingRegionServers         *int64   `json:"numDrainingRegionServers,omitempty"`
	ClusterRequests                  *int64   `json:"clusterRequests,omitempty"`
	ServerCrashTimeNumOps            *int64   `json:"ServerCrashTime_num_ops,omitempty"`
	ServerCrashTimeMin               *int64   `json:"ServerCrashTime_min,omitempty"`
	ServerCrashTimeMax               *int64   `json:"ServerCrashTime_max,omitempty"`
	ServerCrashTimeMean              *int64   `json:"ServerCrashTime_mean,omitempty"`
	ServerCrashTime25ThPercentile    *int64   `json:"ServerCrashTime_25th_percentile,omitempty"`
	ServerCrashTimeMedian            *int64   `json:"ServerCrashTime_median,omitempty"`
	ServerCrashTime75ThPercentile    *int64   `json:"ServerCrashTime_75th_percentile,omitempty"`
	ServerCrashTime90ThPercentile    *int64   `json:"ServerCrashTime_90th_percentile,omitempty"`
	ServerCrashTime95ThPercentile    *int64   `json:"ServerCrashTime_95th_percentile,omitempty"`
	ServerCrashTime98ThPercentile    *int64   `json:"ServerCrashTime_98th_percentile,omitempty"`
	ServerCrashTime99ThPercentile    *int64   `json:"ServerCrashTime_99th_percentile,omitempty"`
	ServerCrashTime999ThPercentile   *int64   `json:"ServerCrashTime_99.9th_percentile,omitempty"`
	ServerCrashSubmittedCount        *int64   `json:"ServerCrashSubmittedCount,omitempty"`
	ServerCrashFailedCount           *int64   `json:"ServerCrashFailedCount,omitempty"`
}
