// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    jVMMetrics, err := UnmarshalJVMMetrics(bytes)
//    bytes, err = jVMMetrics.Marshal()
package hadoop

import "encoding/json"

func UnmarshalJVMMetrics(data []byte) (JVMMetrics, error) {
	var r JVMMetrics
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *JVMMetrics) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type JVMMetrics struct {
	Beans []JvmMetricsBean `json:"beans,omitempty"`
}

type JvmMetricsBean struct {
	Name                       *string     `json:"name,omitempty"`
	ModelerType                *string     `json:"modelerType,omitempty"`
	TagContext                 *string     `json:"tag.Context,omitempty"`
	TagProcessName             *string     `json:"tag.ProcessName,omitempty"`
	TagSessionID               interface{} `json:"tag.SessionId"`
	TagHostname                *string     `json:"tag.Hostname,omitempty"`
	MemNonHeapUsedM            *float64    `json:"MemNonHeapUsedM,omitempty"`
	MemNonHeapCommittedM       *float64    `json:"MemNonHeapCommittedM,omitempty"`
	MemNonHeapMaxM             *int64      `json:"MemNonHeapMaxM,omitempty"`
	MemHeapUsedM               *float64    `json:"MemHeapUsedM,omitempty"`
	MemHeapCommittedM          *float64    `json:"MemHeapCommittedM,omitempty"`
	MemHeapMaxM                *int64      `json:"MemHeapMaxM,omitempty"`
	MemMaxM                    *int64      `json:"MemMaxM,omitempty"`
	GcCount                    *int64      `json:"GcCount,omitempty"`
	GcTimeMillis               *int64      `json:"GcTimeMillis,omitempty"`
	GcNumWarnThresholdExceeded *int64      `json:"GcNumWarnThresholdExceeded,omitempty"`
	GcNumInfoThresholdExceeded *int64      `json:"GcNumInfoThresholdExceeded,omitempty"`
	GcTotalExtraSleepTime      *int64      `json:"GcTotalExtraSleepTime,omitempty"`
	ThreadsNew                 *int64      `json:"ThreadsNew,omitempty"`
	ThreadsRunnable            *int64      `json:"ThreadsRunnable,omitempty"`
	ThreadsBlocked             *int64      `json:"ThreadsBlocked,omitempty"`
	ThreadsWaiting             *int64      `json:"ThreadsWaiting,omitempty"`
	ThreadsTimedWaiting        *int64      `json:"ThreadsTimedWaiting,omitempty"`
	ThreadsTerminated          *int64      `json:"ThreadsTerminated,omitempty"`
	LogFatal                   *int64      `json:"LogFatal,omitempty"`
	LogError                   *int64      `json:"LogError,omitempty"`
	LogWarn                    *int64      `json:"LogWarn,omitempty"`
	LogInfo                    *int64      `json:"LogInfo,omitempty"`
}
