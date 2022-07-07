// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    welcome, err := UnmarshalWelcome(bytes)
//    bytes, err = welcome.Marshal()

package hadoop

import "encoding/json"

func UnmarshalFSNamesystem(data []byte) (FSNamesystem, error) {
	var r FSNamesystem
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FSNamesystem) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FSNamesystem struct {
	Beans []FSNamesystemBean `json:"beans,omitempty"`
}

type FSNamesystemBean struct {
	Name                                         *string `json:"name,omitempty"`
	ModelerType                                  *string `json:"modelerType,omitempty"`
	TagContext                                   *string `json:"tag.Context,omitempty"`
	TagHAState                                   *string `json:"tag.HAState,omitempty"`
	TagTotalSyncTimes                            *string `json:"tag.TotalSyncTimes,omitempty"`
	TagHostname                                  *string `json:"tag.Hostname,omitempty"`
	MissingBlocks                                *int64  `json:"MissingBlocks,omitempty"`
	MissingReplOneBlocks                         *int64  `json:"MissingReplOneBlocks,omitempty"`
	ExpiredHeartbeats                            *int64  `json:"ExpiredHeartbeats,omitempty"`
	TransactionsSinceLastCheckpoint              *int64  `json:"TransactionsSinceLastCheckpoint,omitempty"`
	TransactionsSinceLastLogRoll                 *int64  `json:"TransactionsSinceLastLogRoll,omitempty"`
	LastWrittenTransactionID                     *int64  `json:"LastWrittenTransactionId,omitempty"`
	LastCheckpointTime                           *int64  `json:"LastCheckpointTime,omitempty"`
	CapacityTotal                                *int64  `json:"CapacityTotal,omitempty"`
	CapacityTotalGB                              *int64  `json:"CapacityTotalGB,omitempty"`
	CapacityUsed                                 *int64  `json:"CapacityUsed,omitempty"`
	CapacityUsedGB                               *int64  `json:"CapacityUsedGB,omitempty"`
	CapacityRemaining                            *int64  `json:"CapacityRemaining,omitempty"`
	ProvidedCapacityTotal                        *int64  `json:"ProvidedCapacityTotal,omitempty"`
	CapacityRemainingGB                          *int64  `json:"CapacityRemainingGB,omitempty"`
	CapacityUsedNonDFS                           *int64  `json:"CapacityUsedNonDFS,omitempty"`
	TotalLoad                                    *int64  `json:"TotalLoad,omitempty"`
	SnapshottableDirectories                     *int64  `json:"SnapshottableDirectories,omitempty"`
	Snapshots                                    *int64  `json:"Snapshots,omitempty"`
	NumEncryptionZones                           *int64  `json:"NumEncryptionZones,omitempty"`
	LockQueueLength                              *int64  `json:"LockQueueLength,omitempty"`
	BlocksTotal                                  *int64  `json:"BlocksTotal,omitempty"`
	NumFilesUnderConstruction                    *int64  `json:"NumFilesUnderConstruction,omitempty"`
	NumActiveClients                             *int64  `json:"NumActiveClients,omitempty"`
	FilesTotal                                   *int64  `json:"FilesTotal,omitempty"`
	PendingReplicationBlocks                     *int64  `json:"PendingReplicationBlocks,omitempty"`
	PendingReconstructionBlocks                  *int64  `json:"PendingReconstructionBlocks,omitempty"`
	UnderReplicatedBlocks                        *int64  `json:"UnderReplicatedBlocks,omitempty"`
	LowRedundancyBlocks                          *int64  `json:"LowRedundancyBlocks,omitempty"`
	CorruptBlocks                                *int64  `json:"CorruptBlocks,omitempty"`
	ScheduledReplicationBlocks                   *int64  `json:"ScheduledReplicationBlocks,omitempty"`
	PendingDeletionBlocks                        *int64  `json:"PendingDeletionBlocks,omitempty"`
	LowRedundancyReplicatedBlocks                *int64  `json:"LowRedundancyReplicatedBlocks,omitempty"`
	CorruptReplicatedBlocks                      *int64  `json:"CorruptReplicatedBlocks,omitempty"`
	MissingReplicatedBlocks                      *int64  `json:"MissingReplicatedBlocks,omitempty"`
	MissingReplicationOneBlocks                  *int64  `json:"MissingReplicationOneBlocks,omitempty"`
	HighestPriorityLowRedundancyReplicatedBlocks *int64  `json:"HighestPriorityLowRedundancyReplicatedBlocks,omitempty"`
	HighestPriorityLowRedundancyECBlocks         *int64  `json:"HighestPriorityLowRedundancyECBlocks,omitempty"`
	BytesInFutureReplicatedBlocks                *int64  `json:"BytesInFutureReplicatedBlocks,omitempty"`
	PendingDeletionReplicatedBlocks              *int64  `json:"PendingDeletionReplicatedBlocks,omitempty"`
	TotalReplicatedBlocks                        *int64  `json:"TotalReplicatedBlocks,omitempty"`
	LowRedundancyECBlockGroups                   *int64  `json:"LowRedundancyECBlockGroups,omitempty"`
	CorruptECBlockGroups                         *int64  `json:"CorruptECBlockGroups,omitempty"`
	MissingECBlockGroups                         *int64  `json:"MissingECBlockGroups,omitempty"`
	BytesInFutureECBlockGroups                   *int64  `json:"BytesInFutureECBlockGroups,omitempty"`
	PendingDeletionECBlocks                      *int64  `json:"PendingDeletionECBlocks,omitempty"`
	TotalECBlockGroups                           *int64  `json:"TotalECBlockGroups,omitempty"`
	ExcessBlocks                                 *int64  `json:"ExcessBlocks,omitempty"`
	NumTimedOutPendingReconstructions            *int64  `json:"NumTimedOutPendingReconstructions,omitempty"`
	PostponedMisreplicatedBlocks                 *int64  `json:"PostponedMisreplicatedBlocks,omitempty"`
	PendingDataNodeMessageCount                  *int64  `json:"PendingDataNodeMessageCount,omitempty"`
	MillisSinceLastLoadedEdits                   *int64  `json:"MillisSinceLastLoadedEdits,omitempty"`
	BlockCapacity                                *int64  `json:"BlockCapacity,omitempty"`
	NumLiveDataNodes                             *int64  `json:"NumLiveDataNodes,omitempty"`
	NumDeadDataNodes                             *int64  `json:"NumDeadDataNodes,omitempty"`
	NumDecomLiveDataNodes                        *int64  `json:"NumDecomLiveDataNodes,omitempty"`
	NumDecomDeadDataNodes                        *int64  `json:"NumDecomDeadDataNodes,omitempty"`
	VolumeFailuresTotal                          *int64  `json:"VolumeFailuresTotal,omitempty"`
	EstimatedCapacityLostTotal                   *int64  `json:"EstimatedCapacityLostTotal,omitempty"`
	NumDecommissioningDataNodes                  *int64  `json:"NumDecommissioningDataNodes,omitempty"`
	StaleDataNodes                               *int64  `json:"StaleDataNodes,omitempty"`
	NumStaleStorages                             *int64  `json:"NumStaleStorages,omitempty"`
	TotalSyncCount                               *int64  `json:"TotalSyncCount,omitempty"`
	NumInMaintenanceLiveDataNodes                *int64  `json:"NumInMaintenanceLiveDataNodes,omitempty"`
	NumInMaintenanceDeadDataNodes                *int64  `json:"NumInMaintenanceDeadDataNodes,omitempty"`
	NumEnteringMaintenanceDataNodes              *int64  `json:"NumEnteringMaintenanceDataNodes,omitempty"`
}
