// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    fSNamesystemState, err := UnmarshalFSNamesystemState(bytes)
//    bytes, err = fSNamesystemState.Marshal()
package hadoop

import "encoding/json"

func UnmarshalFSNamesystemState(data []byte) (FSNamesystemState, error) {
	var r FSNamesystemState
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FSNamesystemState) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FSNamesystemState struct {
	Beans []FSNamesystemStateBean `json:"beans,omitempty"`
}

type FSNamesystemStateBean struct {
	Name                            *string `json:"name,omitempty"`
	ModelerType                     *string `json:"modelerType,omitempty"`
	CapacityTotal                   *int64  `json:"CapacityTotal,omitempty"`
	CapacityUsed                    *int64  `json:"CapacityUsed,omitempty"`
	CapacityRemaining               *int64  `json:"CapacityRemaining,omitempty"`
	ProvidedCapacityTotal           *int64  `json:"ProvidedCapacityTotal,omitempty"`
	TotalLoad                       *int64  `json:"TotalLoad,omitempty"`
	SnapshotStats                   *string `json:"SnapshotStats,omitempty"`
	NumEncryptionZones              *int64  `json:"NumEncryptionZones,omitempty"`
	FSLockQueueLength               *int64  `json:"FsLockQueueLength,omitempty"`
	BlocksTotal                     *int64  `json:"BlocksTotal,omitempty"`
	MaxObjects                      *int64  `json:"MaxObjects,omitempty"`
	FilesTotal                      *int64  `json:"FilesTotal,omitempty"`
	PendingReplicationBlocks        *int64  `json:"PendingReplicationBlocks,omitempty"`
	PendingReconstructionBlocks     *int64  `json:"PendingReconstructionBlocks,omitempty"`
	UnderReplicatedBlocks           *int64  `json:"UnderReplicatedBlocks,omitempty"`
	LowRedundancyBlocks             *int64  `json:"LowRedundancyBlocks,omitempty"`
	ScheduledReplicationBlocks      *int64  `json:"ScheduledReplicationBlocks,omitempty"`
	PendingDeletionBlocks           *int64  `json:"PendingDeletionBlocks,omitempty"`
	BlockDeletionStartTime          *int64  `json:"BlockDeletionStartTime,omitempty"`
	FSState                         *string `json:"FSState,omitempty"`
	NumLiveDataNodes                *int64  `json:"NumLiveDataNodes,omitempty"`
	NumDeadDataNodes                *int64  `json:"NumDeadDataNodes,omitempty"`
	NumDecomLiveDataNodes           *int64  `json:"NumDecomLiveDataNodes,omitempty"`
	NumDecomDeadDataNodes           *int64  `json:"NumDecomDeadDataNodes,omitempty"`
	VolumeFailuresTotal             *int64  `json:"VolumeFailuresTotal,omitempty"`
	EstimatedCapacityLostTotal      *int64  `json:"EstimatedCapacityLostTotal,omitempty"`
	NumDecommissioningDataNodes     *int64  `json:"NumDecommissioningDataNodes,omitempty"`
	NumStaleDataNodes               *int64  `json:"NumStaleDataNodes,omitempty"`
	NumStaleStorages                *int64  `json:"NumStaleStorages,omitempty"`
	TopUserOpCounts                 *string `json:"TopUserOpCounts,omitempty"`
	TotalSyncCount                  *int64  `json:"TotalSyncCount,omitempty"`
	TotalSyncTimes                  *string `json:"TotalSyncTimes,omitempty"`
	NumInMaintenanceLiveDataNodes   *int64  `json:"NumInMaintenanceLiveDataNodes,omitempty"`
	NumInMaintenanceDeadDataNodes   *int64  `json:"NumInMaintenanceDeadDataNodes,omitempty"`
	NumEnteringMaintenanceDataNodes *int64  `json:"NumEnteringMaintenanceDataNodes,omitempty"`
}
