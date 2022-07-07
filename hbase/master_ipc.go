// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    masterIPC, err := UnmarshalMasterIPC(bytes)
//    bytes, err = masterIPC.Marshal()

package hbase

import "encoding/json"

func UnmarshalMasterIPC(data []byte) (MasterIPC, error) {
	var r MasterIPC
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MasterIPC) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MasterIPC struct {
	Beans []IPCBean `json:"beans,omitempty"`
}

type IPCBean struct {
	Name                                     *string `json:"name,omitempty"`
	ModelerType                              *string `json:"modelerType,omitempty"`
	TagContext                               *string `json:"tag.Context,omitempty"`
	TagHostname                              *string `json:"tag.Hostname,omitempty"`
	QueueSize                                *int64  `json:"queueSize,omitempty"`
	NumCallsInGeneralQueue                   *int64  `json:"numCallsInGeneralQueue,omitempty"`
	NumCallsInReplicationQueue               *int64  `json:"numCallsInReplicationQueue,omitempty"`
	NumCallsInPriorityQueue                  *int64  `json:"numCallsInPriorityQueue,omitempty"`
	NumCallsInMetaPriorityQueue              *int64  `json:"numCallsInMetaPriorityQueue,omitempty"`
	NumOpenConnections                       *int64  `json:"numOpenConnections,omitempty"`
	NumActiveHandler                         *int64  `json:"numActiveHandler,omitempty"`
	NumActiveGeneralHandler                  *int64  `json:"numActiveGeneralHandler,omitempty"`
	NumActivePriorityHandler                 *int64  `json:"numActivePriorityHandler,omitempty"`
	NumActiveReplicationHandler              *int64  `json:"numActiveReplicationHandler,omitempty"`
	NumGeneralCallsDropped                   *int64  `json:"numGeneralCallsDropped,omitempty"`
	NumLIFOModeSwitches                      *int64  `json:"numLifoModeSwitches,omitempty"`
	NumCallsInWriteQueue                     *int64  `json:"numCallsInWriteQueue,omitempty"`
	NumCallsInReadQueue                      *int64  `json:"numCallsInReadQueue,omitempty"`
	NumCallsInScanQueue                      *int64  `json:"numCallsInScanQueue,omitempty"`
	NumActiveWriteHandler                    *int64  `json:"numActiveWriteHandler,omitempty"`
	NumActiveReadHandler                     *int64  `json:"numActiveReadHandler,omitempty"`
	NumActiveScanHandler                     *int64  `json:"numActiveScanHandler,omitempty"`
	NettyDirectMemoryUsage                   *int64  `json:"nettyDirectMemoryUsage,omitempty"`
	ReceivedBytes                            *int64  `json:"receivedBytes,omitempty"`
	ExceptionsRegionMovedException           *int64  `json:"exceptions.RegionMovedException,omitempty"`
	ExceptionsMultiResponseTooLarge          *int64  `json:"exceptions.multiResponseTooLarge,omitempty"`
	AuthenticationSuccesses                  *int64  `json:"authenticationSuccesses,omitempty"`
	AuthorizationFailures                    *int64  `json:"authorizationFailures,omitempty"`
	TotalCallTimeNumOps                      *int64  `json:"TotalCallTime_num_ops,omitempty"`
	TotalCallTimeMin                         *int64  `json:"TotalCallTime_min,omitempty"`
	TotalCallTimeMax                         *int64  `json:"TotalCallTime_max,omitempty"`
	TotalCallTimeMean                        *int64  `json:"TotalCallTime_mean,omitempty"`
	TotalCallTime25ThPercentile              *int64  `json:"TotalCallTime_25th_percentile,omitempty"`
	TotalCallTimeMedian                      *int64  `json:"TotalCallTime_median,omitempty"`
	TotalCallTime75ThPercentile              *int64  `json:"TotalCallTime_75th_percentile,omitempty"`
	TotalCallTime90ThPercentile              *int64  `json:"TotalCallTime_90th_percentile,omitempty"`
	TotalCallTime95ThPercentile              *int64  `json:"TotalCallTime_95th_percentile,omitempty"`
	TotalCallTime98ThPercentile              *int64  `json:"TotalCallTime_98th_percentile,omitempty"`
	TotalCallTime99ThPercentile              *int64  `json:"TotalCallTime_99th_percentile,omitempty"`
	TotalCallTime999ThPercentile             *int64  `json:"TotalCallTime_99.9th_percentile,omitempty"`
	TotalCallTimeTimeRangeCount01            *int64  `json:"TotalCallTime_TimeRangeCount_0-1,omitempty"`
	ExceptionsRegionTooBusyException         *int64  `json:"exceptions.RegionTooBusyException,omitempty"`
	ExceptionsFailedSanityCheckException     *int64  `json:"exceptions.FailedSanityCheckException,omitempty"`
	ResponseSizeNumOps                       *int64  `json:"ResponseSize_num_ops,omitempty"`
	ResponseSizeMin                          *int64  `json:"ResponseSize_min,omitempty"`
	ResponseSizeMax                          *int64  `json:"ResponseSize_max,omitempty"`
	ResponseSizeMean                         *int64  `json:"ResponseSize_mean,omitempty"`
	ResponseSize25ThPercentile               *int64  `json:"ResponseSize_25th_percentile,omitempty"`
	ResponseSizeMedian                       *int64  `json:"ResponseSize_median,omitempty"`
	ResponseSize75ThPercentile               *int64  `json:"ResponseSize_75th_percentile,omitempty"`
	ResponseSize90ThPercentile               *int64  `json:"ResponseSize_90th_percentile,omitempty"`
	ResponseSize95ThPercentile               *int64  `json:"ResponseSize_95th_percentile,omitempty"`
	ResponseSize98ThPercentile               *int64  `json:"ResponseSize_98th_percentile,omitempty"`
	ResponseSize99ThPercentile               *int64  `json:"ResponseSize_99th_percentile,omitempty"`
	ResponseSize999ThPercentile              *int64  `json:"ResponseSize_99.9th_percentile,omitempty"`
	ResponseSizeSizeRangeCount010            *int64  `json:"ResponseSize_SizeRangeCount_0-10,omitempty"`
	ExceptionsUnknownScannerException        *int64  `json:"exceptions.UnknownScannerException,omitempty"`
	ExceptionsOutOfOrderScannerNextException *int64  `json:"exceptions.OutOfOrderScannerNextException,omitempty"`
	Exceptions                               *int64  `json:"exceptions,omitempty"`
	ProcessCallTimeNumOps                    *int64  `json:"ProcessCallTime_num_ops,omitempty"`
	ProcessCallTimeMin                       *int64  `json:"ProcessCallTime_min,omitempty"`
	ProcessCallTimeMax                       *int64  `json:"ProcessCallTime_max,omitempty"`
	ProcessCallTimeMean                      *int64  `json:"ProcessCallTime_mean,omitempty"`
	ProcessCallTime25ThPercentile            *int64  `json:"ProcessCallTime_25th_percentile,omitempty"`
	ProcessCallTimeMedian                    *int64  `json:"ProcessCallTime_median,omitempty"`
	ProcessCallTime75ThPercentile            *int64  `json:"ProcessCallTime_75th_percentile,omitempty"`
	ProcessCallTime90ThPercentile            *int64  `json:"ProcessCallTime_90th_percentile,omitempty"`
	ProcessCallTime95ThPercentile            *int64  `json:"ProcessCallTime_95th_percentile,omitempty"`
	ProcessCallTime98ThPercentile            *int64  `json:"ProcessCallTime_98th_percentile,omitempty"`
	ProcessCallTime99ThPercentile            *int64  `json:"ProcessCallTime_99th_percentile,omitempty"`
	ProcessCallTime999ThPercentile           *int64  `json:"ProcessCallTime_99.9th_percentile,omitempty"`
	ProcessCallTimeTimeRangeCount01          *int64  `json:"ProcessCallTime_TimeRangeCount_0-1,omitempty"`
	AuthenticationFallbacks                  *int64  `json:"authenticationFallbacks,omitempty"`
	ExceptionsNotServingRegionException      *int64  `json:"exceptions.NotServingRegionException,omitempty"`
	ExceptionsCallQueueTooBig                *int64  `json:"exceptions.callQueueTooBig,omitempty"`
	AuthorizationSuccesses                   *int64  `json:"authorizationSuccesses,omitempty"`
	ExceptionsScannerResetException          *int64  `json:"exceptions.ScannerResetException,omitempty"`
	RequestSizeNumOps                        *int64  `json:"RequestSize_num_ops,omitempty"`
	RequestSizeMin                           *int64  `json:"RequestSize_min,omitempty"`
	RequestSizeMax                           *int64  `json:"RequestSize_max,omitempty"`
	RequestSizeMean                          *int64  `json:"RequestSize_mean,omitempty"`
	RequestSize25ThPercentile                *int64  `json:"RequestSize_25th_percentile,omitempty"`
	RequestSizeMedian                        *int64  `json:"RequestSize_median,omitempty"`
	RequestSize75ThPercentile                *int64  `json:"RequestSize_75th_percentile,omitempty"`
	RequestSize90ThPercentile                *int64  `json:"RequestSize_90th_percentile,omitempty"`
	RequestSize95ThPercentile                *int64  `json:"RequestSize_95th_percentile,omitempty"`
	RequestSize98ThPercentile                *int64  `json:"RequestSize_98th_percentile,omitempty"`
	RequestSize99ThPercentile                *int64  `json:"RequestSize_99th_percentile,omitempty"`
	RequestSize999ThPercentile               *int64  `json:"RequestSize_99.9th_percentile,omitempty"`
	RequestSizeSizeRangeCount1001000         *int64  `json:"RequestSize_SizeRangeCount_100-1000,omitempty"`
	SentBytes                                *int64  `json:"sentBytes,omitempty"`
	QueueCallTimeNumOps                      *int64  `json:"QueueCallTime_num_ops,omitempty"`
	QueueCallTimeMin                         *int64  `json:"QueueCallTime_min,omitempty"`
	QueueCallTimeMax                         *int64  `json:"QueueCallTime_max,omitempty"`
	QueueCallTimeMean                        *int64  `json:"QueueCallTime_mean,omitempty"`
	QueueCallTime25ThPercentile              *int64  `json:"QueueCallTime_25th_percentile,omitempty"`
	QueueCallTimeMedian                      *int64  `json:"QueueCallTime_median,omitempty"`
	QueueCallTime75ThPercentile              *int64  `json:"QueueCallTime_75th_percentile,omitempty"`
	QueueCallTime90ThPercentile              *int64  `json:"QueueCallTime_90th_percentile,omitempty"`
	QueueCallTime95ThPercentile              *int64  `json:"QueueCallTime_95th_percentile,omitempty"`
	QueueCallTime98ThPercentile              *int64  `json:"QueueCallTime_98th_percentile,omitempty"`
	QueueCallTime99ThPercentile              *int64  `json:"QueueCallTime_99th_percentile,omitempty"`
	QueueCallTime999ThPercentile             *int64  `json:"QueueCallTime_99.9th_percentile,omitempty"`
	QueueCallTimeTimeRangeCount01            *int64  `json:"QueueCallTime_TimeRangeCount_0-1,omitempty"`
	AuthenticationFailures                   *int64  `json:"authenticationFailures,omitempty"`
}
