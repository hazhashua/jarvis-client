package micro_service

import "encoding/json"

func UnmarshalK8sNodeInfo(data []byte) (K8sNodeInfo, error) {
	var r K8sNodeInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *K8sNodeInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type K8sNodeInfo struct {
	Kind       *string           `json:"kind,omitempty"`
	APIVersion *string           `json:"apiVersion,omitempty"`
	Metadata   *NodeListMetadata `json:"metadata,omitempty"`
	Items      []NodeItem        `json:"items,omitempty"`
}

type NodeItem struct {
	Metadata *NodeItemMetadata `json:"metadata,omitempty"`
	Spec     *NodeSpec         `json:"spec,omitempty"`
	Status   *StatusClass      `json:"status,omitempty"`
}

type NodeItemMetadata struct {
	Name              *string          `json:"name,omitempty"`
	SelfLink          *string          `json:"selfLink,omitempty"`
	Uid               *string          `json:"uid,omitempty"`
	ResourceVersion   *string          `json:"resourceVersion,omitempty"`
	CreationTimestamp *string          `json:"creationTimestamp,omitempty"`
	Labels            *NodeLabels      `json:"labels,omitempty"`
	Annotations       *NodeAnnotations `json:"annotations,omitempty"`
}

type NodeAnnotations struct {
	NodeAlphaKubernetesIoTTL                         *string `json:"node.alpha.kubernetes.io/ttl,omitempty"`
	VolumesKubernetesIoControllerManagedAttachDetach *string `json:"volumes.kubernetes.io/controller-managed-attach-detach,omitempty"`
}

type NodeLabels struct {
	BetaKubernetesIoArch       *string `json:"beta.kubernetes.io/arch,omitempty"`
	BetaKubernetesIoOS         *string `json:"beta.kubernetes.io/os,omitempty"`
	KubernetesIoArch           *string `json:"kubernetes.io/arch,omitempty"`
	KubernetesIoHostname       *string `json:"kubernetes.io/hostname,omitempty"`
	KubernetesIoOS             *string `json:"kubernetes.io/os,omitempty"`
	NodeRoleKubernetesIoNode   *string `json:"node-role.kubernetes.io/node,omitempty"`
	NodeRoleKubernetesIoMaster *string `json:"node-role.kubernetes.io/master,omitempty"`
	MinioServer                *string `json:"minio-server,omitempty"`
}

type NodeSpec struct {
}

type StatusClass struct {
	Capacity        *Allocatable     `json:"capacity,omitempty"`
	Allocatable     *Allocatable     `json:"allocatable,omitempty"`
	Conditions      []Condition      `json:"conditions,omitempty"`
	Addresses       []NodeAddress    `json:"addresses,omitempty"`
	DaemonEndpoints *DaemonEndpoints `json:"daemonEndpoints,omitempty"`
	NodeInfo        *NodeInfo        `json:"nodeInfo,omitempty"`
	Images          []Image          `json:"images,omitempty"`
}

type NodeAddress struct {
	Type    *AddressType `json:"type,omitempty"`
	Address *string      `json:"address,omitempty"`
}

type Allocatable struct {
	CPU              *string `json:"cpu,omitempty"`
	EphemeralStorage *string `json:"ephemeral-storage,omitempty"`
	Hugepages1Gi     *string `json:"hugepages-1Gi,omitempty"`
	Hugepages2Mi     *string `json:"hugepages-2Mi,omitempty"`
	Memory           *string `json:"memory,omitempty"`
	Pods             *string `json:"pods,omitempty"`
}

type Condition struct {
	Type               *ConditionType `json:"type,omitempty"`
	Status             *StatusEnum    `json:"status,omitempty"`
	LastHeartbeatTime  *string        `json:"lastHeartbeatTime,omitempty"`
	LastTransitionTime *string        `json:"lastTransitionTime,omitempty"`
	Reason             *Reason        `json:"reason,omitempty"`
	Message            *Message       `json:"message,omitempty"`
}

type DaemonEndpoints struct {
	KubeletEndpoint *KubeletEndpoint `json:"kubeletEndpoint,omitempty"`
}

type KubeletEndpoint struct {
	Port *int64 `json:"Port,omitempty"`
}

type Image struct {
	Names     []string `json:"names,omitempty"`
	SizeBytes *int64   `json:"sizeBytes,omitempty"`
}

type NodeInfo struct {
	MachineID               *string `json:"machineID,omitempty"`
	SystemUUID              *string `json:"systemUUID,omitempty"`
	BootID                  *string `json:"bootID,omitempty"`
	KernelVersion           *string `json:"kernelVersion,omitempty"`
	OSImage                 *string `json:"osImage,omitempty"`
	ContainerRuntimeVersion *string `json:"containerRuntimeVersion,omitempty"`
	KubeletVersion          *string `json:"kubeletVersion,omitempty"`
	KubeProxyVersion        *string `json:"kubeProxyVersion,omitempty"`
	OperatingSystem         *string `json:"operatingSystem,omitempty"`
	Architecture            *string `json:"architecture,omitempty"`
}

type NodeListMetadata struct {
	SelfLink        *string `json:"selfLink,omitempty"`
	ResourceVersion *string `json:"resourceVersion,omitempty"`
}

type AddressType string

const (
	Hostname   AddressType = "Hostname"
	InternalIP AddressType = "InternalIP"
)

type Message string

const (
	KubeletHasNoDiskPressure            Message = "kubelet has no disk pressure"
	KubeletHasSufficientMemoryAvailable Message = "kubelet has sufficient memory available"
	KubeletHasSufficientPIDAvailable    Message = "kubelet has sufficient PID available"
	KubeletIsPostingReadyStatus         Message = "kubelet is posting ready status"
)

type Reason string

const (
	KubeletHasSufficientMemory     Reason = "KubeletHasSufficientMemory"
	KubeletHasSufficientPID        Reason = "KubeletHasSufficientPID"
	KubeletReady                   Reason = "KubeletReady"
	ReasonKubeletHasNoDiskPressure Reason = "KubeletHasNoDiskPressure"
)

type StatusEnum string

const (
	False StatusEnum = "False"
	True  StatusEnum = "True"
)

type ConditionType string

const (
	DiskPressure   ConditionType = "DiskPressure"
	MemoryPressure ConditionType = "MemoryPressure"
	PIDPressure    ConditionType = "PIDPressure"
	Ready          ConditionType = "Ready"
)
