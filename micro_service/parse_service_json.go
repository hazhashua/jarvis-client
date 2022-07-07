package micro_service

import "encoding/json"

func UnmarshalAPIV1Services(data []byte) (APIV1Services, error) {
	var r APIV1Services
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *APIV1Services) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type APIV1Services struct {
	Kind       *string                `json:"kind,omitempty"`
	APIVersion *string                `json:"apiVersion,omitempty"`
	Metadata   *APIV1ServicesMetadata `json:"metadata,omitempty"`
	Items      []Item                 `json:"items,omitempty"`
}

type Item struct {
	Metadata *ItemMetadata `json:"metadata,omitempty"`
	Spec     *Spec         `json:"spec,omitempty"`
	Status   *Status       `json:"status,omitempty"`
}

type ItemMetadata struct {
	Name              *string      `json:"name,omitempty"`
	Namespace         *Namespace   `json:"namespace,omitempty"`
	SelfLink          *string      `json:"selfLink,omitempty"`
	Uid               *string      `json:"uid,omitempty"`
	ResourceVersion   *string      `json:"resourceVersion,omitempty"`
	CreationTimestamp *string      `json:"creationTimestamp,omitempty"`
	Labels            *Labels      `json:"labels,omitempty"`
	Annotations       *Annotations `json:"annotations,omitempty"`
}

type Annotations struct {
	KubectlKubernetesIoLastAppliedConfiguration *string `json:"kubectl.kubernetes.io/last-applied-configuration,omitempty"`
	PrometheusIoPort                            *string `json:"prometheus.io/port,omitempty"`
	PrometheusIoScrape                          *string `json:"prometheus.io/scrape,omitempty"`
}

type Labels struct {
	App                          *string `json:"app,omitempty"`
	Component                    *string `json:"component,omitempty"`
	Provider                     *string `json:"provider,omitempty"`
	AddonmanagerKubernetesIoMode *string `json:"addonmanager.kubernetes.io/mode,omitempty"`
	K8SApp                       *string `json:"k8s-app,omitempty"`
	KubernetesIoClusterService   *string `json:"kubernetes.io/cluster-service,omitempty"`
	KubernetesIoName             *string `json:"kubernetes.io/name,omitempty"`
}

type Spec struct {
	Ports                 []Port                 `json:"ports,omitempty"`
	Selector              *Selector              `json:"selector,omitempty"`
	ClusterIP             *string                `json:"clusterIP,omitempty"`
	Type                  *Type                  `json:"type,omitempty"`
	SessionAffinity       *SessionAffinity       `json:"sessionAffinity,omitempty"`
	ExternalTrafficPolicy *ExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty"`
	SessionAffinityConfig *SessionAffinityConfig `json:"sessionAffinityConfig,omitempty"`
}

type Port struct {
	Name       *string   `json:"name,omitempty"`
	Protocol   *Protocol `json:"protocol,omitempty"`
	Port       *int64    `json:"port,omitempty"`
	TargetPort *int64    `json:"targetPort,omitempty"`
	NodePort   *int64    `json:"nodePort,omitempty"`
}

type Selector struct {
	App    *string `json:"app,omitempty"`
	K8SApp *string `json:"k8s-app,omitempty"`
}

type SessionAffinityConfig struct {
	ClientIP *ClientIPClass `json:"clientIP,omitempty"`
}

type ClientIPClass struct {
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`
}

type Status struct {
	LoadBalancer *LoadBalancer `json:"loadBalancer,omitempty"`
}

type LoadBalancer struct {
}

type APIV1ServicesMetadata struct {
	SelfLink        *string `json:"selfLink,omitempty"`
	ResourceVersion *string `json:"resourceVersion,omitempty"`
}

type Namespace string

const (
	Default    Namespace = "default"
	KubeSystem Namespace = "kube-system"
)

type ExternalTrafficPolicy string

const (
	Cluster ExternalTrafficPolicy = "Cluster"
	Local   ExternalTrafficPolicy = "Local"
)

type Protocol string

const (
	TCP Protocol = "TCP"
	UDP Protocol = "UDP"
)

type SessionAffinity string

const (
	ClientIP SessionAffinity = "ClientIP"
	None     SessionAffinity = "None"
)

type Type string

const (
	ClusterIP Type = "ClusterIP"
	NodePort  Type = "NodePort"
)
