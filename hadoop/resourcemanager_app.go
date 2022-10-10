// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    resourceManagerApp, err := UnmarshalResourceManagerApp(bytes)
//    bytes, err = resourceManagerApp.Marshal()
package hadoop

import "encoding/json"

func UnmarshalResourceManagerApp(data []byte) (ResourceManagerApp, error) {
	var r ResourceManagerApp
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ResourceManagerApp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ResourceManagerApp struct {
	Beans []QueueBean `json:"beans,omitempty"`
}

type QueueBean struct {
	Name                                           *string  `json:"name,omitempty"`
	ModelerType                                    *string  `json:"modelerType,omitempty"`
	TagQueue                                       *string  `json:"tag.Queue,omitempty"`
	TagContext                                     *string  `json:"tag.Context,omitempty"`
	TagHostname                                    *string  `json:"tag.Hostname,omitempty"`
	Running0                                       *int64   `json:"running_0,omitempty"`
	Running60                                      *int64   `json:"running_60,omitempty"`
	Running300                                     *int64   `json:"running_300,omitempty"`
	Running1440                                    *int64   `json:"running_1440,omitempty"`
	AMResourceLimitMB                              *int64   `json:"AMResourceLimitMB,omitempty"`
	AMResourceLimitVCores                          *int64   `json:"AMResourceLimitVCores,omitempty"`
	UsedAMResourceMB                               *int64   `json:"UsedAMResourceMB,omitempty"`
	UsedAMResourceVCores                           *int64   `json:"UsedAMResourceVCores,omitempty"`
	UsedCapacity                                   *float32 `json:"UsedCapacity,omitempty"`
	AbsoluteUsedCapacity                           *float32 `json:"AbsoluteUsedCapacity,omitempty"`
	GuaranteedMB                                   *int64   `json:"GuaranteedMB,omitempty"`
	GuaranteedVCores                               *int64   `json:"GuaranteedVCores,omitempty"`
	MaxCapacityMB                                  *int64   `json:"MaxCapacityMB,omitempty"`
	MaxCapacityVCores                              *int64   `json:"MaxCapacityVCores,omitempty"`
	AppsSubmitted                                  *int64   `json:"AppsSubmitted,omitempty"`
	AppsRunning                                    *int64   `json:"AppsRunning,omitempty"`
	AppsPending                                    *int64   `json:"AppsPending,omitempty"`
	AppsCompleted                                  *int64   `json:"AppsCompleted,omitempty"`
	AppsKilled                                     *int64   `json:"AppsKilled,omitempty"`
	AppsFailed                                     *int64   `json:"AppsFailed,omitempty"`
	AggregateNodeLocalContainersAllocated          *int64   `json:"AggregateNodeLocalContainersAllocated,omitempty"`
	AggregateRackLocalContainersAllocated          *int64   `json:"AggregateRackLocalContainersAllocated,omitempty"`
	AggregateOffSwitchContainersAllocated          *int64   `json:"AggregateOffSwitchContainersAllocated,omitempty"`
	AggregateContainersPreempted                   *int64   `json:"AggregateContainersPreempted,omitempty"`
	AggregateMemoryMBSecondsPreempted              *int64   `json:"AggregateMemoryMBSecondsPreempted,omitempty"`
	AggregateVcoreSecondsPreempted                 *int64   `json:"AggregateVcoreSecondsPreempted,omitempty"`
	ActiveUsers                                    *int64   `json:"ActiveUsers,omitempty"`
	ActiveApplications                             *int64   `json:"ActiveApplications,omitempty"`
	AppAttemptFirstContainerAllocationDelayNumOps  *int64   `json:"AppAttemptFirstContainerAllocationDelayNumOps,omitempty"`
	AppAttemptFirstContainerAllocationDelayAvgTime *float32 `json:"AppAttemptFirstContainerAllocationDelayAvgTime,omitempty"`
	AllocatedMB                                    *int64   `json:"AllocatedMB,omitempty"`
	AllocatedVCores                                *int64   `json:"AllocatedVCores,omitempty"`
	AllocatedContainers                            *int64   `json:"AllocatedContainers,omitempty"`
	AggregateContainersAllocated                   *int64   `json:"AggregateContainersAllocated,omitempty"`
	AggregateContainersReleased                    *int64   `json:"AggregateContainersReleased,omitempty"`
	AvailableMB                                    *int64   `json:"AvailableMB,omitempty"`
	AvailableVCores                                *int64   `json:"AvailableVCores,omitempty"`
	PendingMB                                      *int64   `json:"PendingMB,omitempty"`
	PendingVCores                                  *int64   `json:"PendingVCores,omitempty"`
	PendingContainers                              *int64   `json:"PendingContainers,omitempty"`
	ReservedMB                                     *int64   `json:"ReservedMB,omitempty"`
	ReservedVCores                                 *int64   `json:"ReservedVCores,omitempty"`
	ReservedContainers                             *int64   `json:"ReservedContainers,omitempty"`
}
