package micro_service

import (
	"encoding/json"
)

func UnmarshalAPIV1Endpoints(data []byte) (APIV1Endpoints, error) {
	var r APIV1Endpoints
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *APIV1Endpoints) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type APIV1Endpoints struct {
	Kind       *string                 `json:"kind,omitempty"`
	APIVersion *string                 `json:"apiVersion,omitempty"`
	Metadata   *APIV1EndpointsMetadata `json:"metadata,omitempty"`
	Items      []EndpointItem          `json:"items,omitempty"`
}

type EndpointItem struct {
	Metadata *EndpointItemMetadata `json:"metadata,omitempty"`
	Subsets  []Subset              `json:"subsets,omitempty"`
}

type EndpointItemMetadata struct {
	Name              *string              `json:"name,omitempty"`
	Namespace         *Namespace           `json:"namespace,omitempty"`
	SelfLink          *string              `json:"selfLink,omitempty"`
	Uid               *string              `json:"uid,omitempty"`
	ResourceVersion   *string              `json:"resourceVersion,omitempty"`
	CreationTimestamp *string              `json:"creationTimestamp,omitempty"`
	Labels            *EndpointLabels      `json:"labels,omitempty"`
	Annotations       *EndpointAnnotations `json:"annotations,omitempty"`
}

type EndpointAnnotations struct {
	EndpointsKubernetesIoLastChangeTriggerTime *string `json:"endpoints.kubernetes.io/last-change-trigger-time,omitempty"`
	ControlPlaneAlphaKubernetesIoLeader        *string `json:"control-plane.alpha.kubernetes.io/leader,omitempty"`
}

type EndpointLabels struct {
	App                          *string `json:"app,omitempty"`
	AddonmanagerKubernetesIoMode *string `json:"addonmanager.kubernetes.io/mode,omitempty"`
	K8SApp                       *string `json:"k8s-app,omitempty"`
	KubernetesIoClusterService   *string `json:"kubernetes.io/cluster-service,omitempty"`
	KubernetesIoName             *string `json:"kubernetes.io/name,omitempty"`
}

type Subset struct {
	Addresses []Address      `json:"addresses,omitempty"`
	Ports     []EndpointPort `json:"ports,omitempty"`
}

type Address struct {
	IP        *string    `json:"ip,omitempty"`
	NodeName  *string    `json:"nodeName,omitempty"`
	TargetRef *TargetRef `json:"targetRef,omitempty"`
	Hostname  *string    `json:"hostname,omitempty"`
}

type TargetRef struct {
	Kind            *Kind      `json:"kind,omitempty"`
	Namespace       *Namespace `json:"namespace,omitempty"`
	Name            *string    `json:"name,omitempty"`
	Uid             *string    `json:"uid,omitempty"`
	ResourceVersion *string    `json:"resourceVersion,omitempty"`
}

type EndpointPort struct {
	Name     *string   `json:"name,omitempty"`
	Port     *int64    `json:"port,omitempty"`
	Protocol *Protocol `json:"protocol,omitempty"`
}

type APIV1EndpointsMetadata struct {
	SelfLink        *string `json:"selfLink,omitempty"`
	ResourceVersion *string `json:"resourceVersion,omitempty"`
}

// type Namespace string
// const (
// 	Default    Namespace = "default"
// 	KubeSystem Namespace = "kube-system"
// )

type Kind string

const (
	Pod Kind = "Pod"
)

// type Protocol string
// const (
// 	TCP Protocol = "TCP"
// 	UDP Protocol = "UDP"
// )
