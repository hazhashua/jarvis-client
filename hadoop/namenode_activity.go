// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    nameNodeActivity, err := UnmarshalNameNodeActivity(bytes)
//    bytes, err = nameNodeActivity.Marshal()
package hadoop

import "encoding/json"

func UnmarshalNameNodeActivity(data []byte) (NameNodeActivity, error) {
	var r NameNodeActivity
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NameNodeActivity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NameNodeActivity struct {
	Beans []NameNodeActivityBean `json:"beans,omitempty"`
}

type NameNodeActivityBean struct {
	Name                              *string     `json:"name,omitempty"`
	ModelerType                       *string     `json:"modelerType,omitempty"`
	TagProcessName                    *string     `json:"tag.ProcessName,omitempty"`
	TagSessionID                      interface{} `json:"tag.SessionId"`
	TagContext                        *string     `json:"tag.Context,omitempty"`
	TagHostname                       *string     `json:"tag.Hostname,omitempty"`
	CreateFileOps                     *int64      `json:"CreateFileOps,omitempty"`
	FilesCreated                      *int64      `json:"FilesCreated,omitempty"`
	FilesAppended                     *int64      `json:"FilesAppended,omitempty"`
	GetBlockLocations                 *int64      `json:"GetBlockLocations,omitempty"`
	FilesRenamed                      *int64      `json:"FilesRenamed,omitempty"`
	FilesTruncated                    *int64      `json:"FilesTruncated,omitempty"`
	GetListingOps                     *int64      `json:"GetListingOps,omitempty"`
	DeleteFileOps                     *int64      `json:"DeleteFileOps,omitempty"`
	FilesDeleted                      *int64      `json:"FilesDeleted,omitempty"`
	FileInfoOps                       *int64      `json:"FileInfoOps,omitempty"`
	AddBlockOps                       *int64      `json:"AddBlockOps,omitempty"`
	GetAdditionalDatanodeOps          *int64      `json:"GetAdditionalDatanodeOps,omitempty"`
	CreateSymlinkOps                  *int64      `json:"CreateSymlinkOps,omitempty"`
	GetLinkTargetOps                  *int64      `json:"GetLinkTargetOps,omitempty"`
	FilesInGetListingOps              *int64      `json:"FilesInGetListingOps,omitempty"`
	SuccessfulReReplications          *int64      `json:"SuccessfulReReplications,omitempty"`
	NumTimesReReplicationNotScheduled *int64      `json:"NumTimesReReplicationNotScheduled,omitempty"`
	TimeoutReReplications             *int64      `json:"TimeoutReReplications,omitempty"`
	AllowSnapshotOps                  *int64      `json:"AllowSnapshotOps,omitempty"`
	DisallowSnapshotOps               *int64      `json:"DisallowSnapshotOps,omitempty"`
	CreateSnapshotOps                 *int64      `json:"CreateSnapshotOps,omitempty"`
	DeleteSnapshotOps                 *int64      `json:"DeleteSnapshotOps,omitempty"`
	RenameSnapshotOps                 *int64      `json:"RenameSnapshotOps,omitempty"`
	ListSnapshottableDirOps           *int64      `json:"ListSnapshottableDirOps,omitempty"`
	SnapshotDiffReportOps             *int64      `json:"SnapshotDiffReportOps,omitempty"`
	BlockReceivedAndDeletedOps        *int64      `json:"BlockReceivedAndDeletedOps,omitempty"`
	BlockOpsQueued                    *int64      `json:"BlockOpsQueued,omitempty"`
	BlockOpsBatched                   *int64      `json:"BlockOpsBatched,omitempty"`
	TransactionsNumOps                *int64      `json:"TransactionsNumOps,omitempty"`
	TransactionsAvgTime               *float64    `json:"TransactionsAvgTime,omitempty"`
	SyncsNumOps                       *int64      `json:"SyncsNumOps,omitempty"`
	SyncsAvgTime                      *float64    `json:"SyncsAvgTime,omitempty"`
	TransactionsBatchedInSync         *int64      `json:"TransactionsBatchedInSync,omitempty"`
	StorageBlockReportNumOps          *int64      `json:"StorageBlockReportNumOps,omitempty"`
	StorageBlockReportAvgTime         *int64      `json:"StorageBlockReportAvgTime,omitempty"`
	CacheReportNumOps                 *int64      `json:"CacheReportNumOps,omitempty"`
	CacheReportAvgTime                *int64      `json:"CacheReportAvgTime,omitempty"`
	GenerateEDEKTimeNumOps            *int64      `json:"GenerateEDEKTimeNumOps,omitempty"`
	GenerateEDEKTimeAvgTime           *int64      `json:"GenerateEDEKTimeAvgTime,omitempty"`
	WarmUpEDEKTimeNumOps              *int64      `json:"WarmUpEDEKTimeNumOps,omitempty"`
	WarmUpEDEKTimeAvgTime             *int64      `json:"WarmUpEDEKTimeAvgTime,omitempty"`
	ResourceCheckTimeNumOps           *int64      `json:"ResourceCheckTimeNumOps,omitempty"`
	ResourceCheckTimeAvgTime          *float64    `json:"ResourceCheckTimeAvgTime,omitempty"`
	SafeModeTime                      *int64      `json:"SafeModeTime,omitempty"`
	FSImageLoadTime                   *int64      `json:"FsImageLoadTime,omitempty"`
	EditLogTailTimeNumOps             *int64      `json:"EditLogTailTimeNumOps,omitempty"`
	EditLogTailTimeAvgTime            *int64      `json:"EditLogTailTimeAvgTime,omitempty"`
	EditLogFetchTimeNumOps            *int64      `json:"EditLogFetchTimeNumOps,omitempty"`
	EditLogFetchTimeAvgTime           *float64    `json:"EditLogFetchTimeAvgTime,omitempty"`
	NumEditLogLoadedNumOps            *int64      `json:"NumEditLogLoadedNumOps,omitempty"`
	NumEditLogLoadedAvgCount          *int64      `json:"NumEditLogLoadedAvgCount,omitempty"`
	EditLogTailIntervalNumOps         *int64      `json:"EditLogTailIntervalNumOps,omitempty"`
	EditLogTailIntervalAvgTime        *int64      `json:"EditLogTailIntervalAvgTime,omitempty"`
	GetEditNumOps                     *int64      `json:"GetEditNumOps,omitempty"`
	GetEditAvgTime                    *int64      `json:"GetEditAvgTime,omitempty"`
	GetImageNumOps                    *int64      `json:"GetImageNumOps,omitempty"`
	GetImageAvgTime                   *int64      `json:"GetImageAvgTime,omitempty"`
	PutImageNumOps                    *int64      `json:"PutImageNumOps,omitempty"`
	PutImageAvgTime                   *int64      `json:"PutImageAvgTime,omitempty"`
	TotalFileOps                      *int64      `json:"TotalFileOps,omitempty"`
}
