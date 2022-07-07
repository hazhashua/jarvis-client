// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    dataNodeActivity, err := UnmarshalDataNodeActivity(bytes)
//    bytes, err = dataNodeActivity.Marshal()

package hadoop

import "encoding/json"

func UnmarshalDataNodeActivity(data []byte) (DataNodeActivity, error) {
	var r DataNodeActivity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DataNodeActivity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DataNodeActivity struct {
	Beans []DataNodeActivityBean `json:"beans,omitempty"`
}

type DataNodeActivityBean struct {
	Name                                       *string     `json:"name,omitempty"`
	ModelerType                                *string     `json:"modelerType,omitempty"`
	TagSessionID                               interface{} `json:"tag.SessionId"`
	TagContext                                 *string     `json:"tag.Context,omitempty"`
	TagHostname                                *string     `json:"tag.Hostname,omitempty"`
	BytesWritten                               *int64      `json:"BytesWritten,omitempty"`
	TotalWriteTime                             *int64      `json:"TotalWriteTime,omitempty"`
	BytesRead                                  *int64      `json:"BytesRead,omitempty"`
	TotalReadTime                              *int64      `json:"TotalReadTime,omitempty"`
	BlocksWritten                              *int64      `json:"BlocksWritten,omitempty"`
	BlocksRead                                 *int64      `json:"BlocksRead,omitempty"`
	BlocksReplicated                           *int64      `json:"BlocksReplicated,omitempty"`
	BlocksRemoved                              *int64      `json:"BlocksRemoved,omitempty"`
	BlocksVerified                             *int64      `json:"BlocksVerified,omitempty"`
	BlockVerificationFailures                  *int64      `json:"BlockVerificationFailures,omitempty"`
	BlocksCached                               *int64      `json:"BlocksCached,omitempty"`
	BlocksUncached                             *int64      `json:"BlocksUncached,omitempty"`
	ReadsFromLocalClient                       *int64      `json:"ReadsFromLocalClient,omitempty"`
	ReadsFromRemoteClient                      *int64      `json:"ReadsFromRemoteClient,omitempty"`
	WritesFromLocalClient                      *int64      `json:"WritesFromLocalClient,omitempty"`
	WritesFromRemoteClient                     *int64      `json:"WritesFromRemoteClient,omitempty"`
	BlocksGetLocalPathInfo                     *int64      `json:"BlocksGetLocalPathInfo,omitempty"`
	RemoteBytesRead                            *int64      `json:"RemoteBytesRead,omitempty"`
	RemoteBytesWritten                         *int64      `json:"RemoteBytesWritten,omitempty"`
	RAMDiskBlocksWrite                         *int64      `json:"RamDiskBlocksWrite,omitempty"`
	RAMDiskBlocksWriteFallback                 *int64      `json:"RamDiskBlocksWriteFallback,omitempty"`
	RAMDiskBytesWrite                          *int64      `json:"RamDiskBytesWrite,omitempty"`
	RAMDiskBlocksReadHits                      *int64      `json:"RamDiskBlocksReadHits,omitempty"`
	RAMDiskBlocksEvicted                       *int64      `json:"RamDiskBlocksEvicted,omitempty"`
	RAMDiskBlocksEvictedWithoutRead            *int64      `json:"RamDiskBlocksEvictedWithoutRead,omitempty"`
	RAMDiskBlocksEvictionWindowMSNumOps        *int64      `json:"RamDiskBlocksEvictionWindowMsNumOps,omitempty"`
	RAMDiskBlocksEvictionWindowMSAvgTime       *int64      `json:"RamDiskBlocksEvictionWindowMsAvgTime,omitempty"`
	RAMDiskBlocksLazyPersisted                 *int64      `json:"RamDiskBlocksLazyPersisted,omitempty"`
	RAMDiskBlocksDeletedBeforeLazyPersisted    *int64      `json:"RamDiskBlocksDeletedBeforeLazyPersisted,omitempty"`
	RAMDiskBytesLazyPersisted                  *int64      `json:"RamDiskBytesLazyPersisted,omitempty"`
	RAMDiskBlocksLazyPersistWindowMSNumOps     *int64      `json:"RamDiskBlocksLazyPersistWindowMsNumOps,omitempty"`
	RAMDiskBlocksLazyPersistWindowMSAvgTime    *int64      `json:"RamDiskBlocksLazyPersistWindowMsAvgTime,omitempty"`
	FsyncCount                                 *int64      `json:"FsyncCount,omitempty"`
	VolumeFailures                             *int64      `json:"VolumeFailures,omitempty"`
	DatanodeNetworkErrors                      *int64      `json:"DatanodeNetworkErrors,omitempty"`
	DataNodeActiveXceiversCount                *int64      `json:"DataNodeActiveXceiversCount,omitempty"`
	ReadBlockOpNumOps                          *int64      `json:"ReadBlockOpNumOps,omitempty"`
	ReadBlockOpAvgTime                         *float64    `json:"ReadBlockOpAvgTime,omitempty"`
	WriteBlockOpNumOps                         *int64      `json:"WriteBlockOpNumOps,omitempty"`
	WriteBlockOpAvgTime                        *float64    `json:"WriteBlockOpAvgTime,omitempty"`
	BlockChecksumOpNumOps                      *int64      `json:"BlockChecksumOpNumOps,omitempty"`
	BlockChecksumOpAvgTime                     *int64      `json:"BlockChecksumOpAvgTime,omitempty"`
	CopyBlockOpNumOps                          *int64      `json:"CopyBlockOpNumOps,omitempty"`
	CopyBlockOpAvgTime                         *int64      `json:"CopyBlockOpAvgTime,omitempty"`
	ReplaceBlockOpNumOps                       *int64      `json:"ReplaceBlockOpNumOps,omitempty"`
	ReplaceBlockOpAvgTime                      *int64      `json:"ReplaceBlockOpAvgTime,omitempty"`
	HeartbeatsNumOps                           *int64      `json:"HeartbeatsNumOps,omitempty"`
	HeartbeatsAvgTime                          *float64    `json:"HeartbeatsAvgTime,omitempty"`
	HeartbeatsTotalNumOps                      *int64      `json:"HeartbeatsTotalNumOps,omitempty"`
	HeartbeatsTotalAvgTime                     *float64    `json:"HeartbeatsTotalAvgTime,omitempty"`
	LifelinesNumOps                            *int64      `json:"LifelinesNumOps,omitempty"`
	LifelinesAvgTime                           *int64      `json:"LifelinesAvgTime,omitempty"`
	BlockReportsNumOps                         *int64      `json:"BlockReportsNumOps,omitempty"`
	BlockReportsAvgTime                        *int64      `json:"BlockReportsAvgTime,omitempty"`
	IncrementalBlockReportsNumOps              *int64      `json:"IncrementalBlockReportsNumOps,omitempty"`
	IncrementalBlockReportsAvgTime             *float64    `json:"IncrementalBlockReportsAvgTime,omitempty"`
	CacheReportsNumOps                         *int64      `json:"CacheReportsNumOps,omitempty"`
	CacheReportsAvgTime                        *int64      `json:"CacheReportsAvgTime,omitempty"`
	PacketACKRoundTripTimeNanosNumOps          *int64      `json:"PacketAckRoundTripTimeNanosNumOps,omitempty"`
	PacketACKRoundTripTimeNanosAvgTime         *float64    `json:"PacketAckRoundTripTimeNanosAvgTime,omitempty"`
	FlushNanosNumOps                           *int64      `json:"FlushNanosNumOps,omitempty"`
	FlushNanosAvgTime                          *float64    `json:"FlushNanosAvgTime,omitempty"`
	FsyncNanosNumOps                           *int64      `json:"FsyncNanosNumOps,omitempty"`
	FsyncNanosAvgTime                          *int64      `json:"FsyncNanosAvgTime,omitempty"`
	SendDataPacketBlockedOnNetworkNanosNumOps  *int64      `json:"SendDataPacketBlockedOnNetworkNanosNumOps,omitempty"`
	SendDataPacketBlockedOnNetworkNanosAvgTime *float64    `json:"SendDataPacketBlockedOnNetworkNanosAvgTime,omitempty"`
	SendDataPacketTransferNanosNumOps          *int64      `json:"SendDataPacketTransferNanosNumOps,omitempty"`
	SendDataPacketTransferNanosAvgTime         *float64    `json:"SendDataPacketTransferNanosAvgTime,omitempty"`
	BlocksInPendingIBR                         *int64      `json:"BlocksInPendingIBR,omitempty"`
	BlocksReceivingInPendingIBR                *int64      `json:"BlocksReceivingInPendingIBR,omitempty"`
	BlocksReceivedInPendingIBR                 *int64      `json:"BlocksReceivedInPendingIBR,omitempty"`
	BlocksDeletedInPendingIBR                  *int64      `json:"BlocksDeletedInPendingIBR,omitempty"`
	EcReconstructionTasks                      *int64      `json:"EcReconstructionTasks,omitempty"`
	EcFailedReconstructionTasks                *int64      `json:"EcFailedReconstructionTasks,omitempty"`
	EcDecodingTimeNanos                        *int64      `json:"EcDecodingTimeNanos,omitempty"`
	EcReconstructionBytesRead                  *int64      `json:"EcReconstructionBytesRead,omitempty"`
	EcReconstructionBytesWritten               *int64      `json:"EcReconstructionBytesWritten,omitempty"`
	EcReconstructionRemoteBytesRead            *int64      `json:"EcReconstructionRemoteBytesRead,omitempty"`
	EcReconstructionReadTimeMillis             *int64      `json:"EcReconstructionReadTimeMillis,omitempty"`
	EcReconstructionDecodingTimeMillis         *int64      `json:"EcReconstructionDecodingTimeMillis,omitempty"`
	EcReconstructionWriteTimeMillis            *int64      `json:"EcReconstructionWriteTimeMillis,omitempty"`
}
