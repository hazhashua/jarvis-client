package micro_service

import "encoding/json"

func UnmarshalPodInfo(data []byte) (PodInfo, error) {
	var r PodInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PodInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PodInfo struct {
	Kind       *string          `json:"kind,omitempty"`
	APIVersion *string          `json:"apiVersion,omitempty"`
	Metadata   *PodInfoMetadata `json:"metadata,omitempty"`
	Items      []PodInfoItem    `json:"items,omitempty"`
}

type PodInfoItem struct {
	Metadata *PodItemMetadata `json:"metadata,omitempty"`
	Spec     *PodSpec         `json:"spec,omitempty"`
	Status   *PodStatusClass  `json:"status,omitempty"`
}

type PodItemMetadata struct {
	Name              *string          `json:"name,omitempty"`
	GenerateName      *string          `json:"generateName,omitempty"`
	Namespace         *Namespace       `json:"namespace,omitempty"`
	SelfLink          *string          `json:"selfLink,omitempty"`
	Uid               *string          `json:"uid,omitempty"`
	ResourceVersion   *string          `json:"resourceVersion,omitempty"`
	CreationTimestamp *string          `json:"creationTimestamp,omitempty"`
	Labels            *PodLabels       `json:"labels,omitempty"`
	OwnerReferences   []OwnerReference `json:"ownerReferences,omitempty"`
	Annotations       *PodAnnotations  `json:"annotations,omitempty"`
}

type PodAnnotations struct {
	SeccompSecurityAlphaKubernetesIoPod *string `json:"seccomp.security.alpha.kubernetes.io/pod,omitempty"`
	PrometheusIoPort                    *string `json:"prometheus.io/port,omitempty"`
	PrometheusIoScrape                  *string `json:"prometheus.io/scrape,omitempty"`
}

type PodLabels struct {
	App                      *string         `json:"app,omitempty"`
	PodTemplateHash          *string         `json:"pod-template-hash,omitempty"`
	ControllerRevisionHash   *string         `json:"controller-revision-hash,omitempty"`
	PodTemplateGeneration    *string         `json:"pod-template-generation,omitempty"`
	ControllerUid            *string         `json:"controller-uid,omitempty"`
	JobName                  *string         `json:"job-name,omitempty"`
	Jobgroup                 *string         `json:"jobgroup,omitempty"`
	Run                      *string         `json:"run,omitempty"`
	K8SApp                   *string         `json:"k8s-app,omitempty"`
	AppKubernetesIoComponent *string         `json:"app.kubernetes.io/component,omitempty"`
	AppKubernetesIoName      *ServiceAccount `json:"app.kubernetes.io/name,omitempty"`
	AppKubernetesIoVersion   *string         `json:"app.kubernetes.io/version,omitempty"`
}

type OwnerReference struct {
	APIVersion         *APIVersion `json:"apiVersion,omitempty"`
	Kind               *Kind       `json:"kind,omitempty"`
	Name               *string     `json:"name,omitempty"`
	Uid                *string     `json:"uid,omitempty"`
	Controller         *bool       `json:"controller,omitempty"`
	BlockOwnerDeletion *bool       `json:"blockOwnerDeletion,omitempty"`
}

type PodSpec struct {
	Volumes                       []Volume        `json:"volumes,omitempty"`
	InitContainers                []InitContainer `json:"initContainers,omitempty"`
	Containers                    []Container     `json:"containers,omitempty"`
	RestartPolicy                 *RestartPolicy  `json:"restartPolicy,omitempty"`
	TerminationGracePeriodSeconds *int64          `json:"terminationGracePeriodSeconds,omitempty"`
	DNSPolicy                     *DNSPolicy      `json:"dnsPolicy,omitempty"`
	ServiceAccountName            *ServiceAccount `json:"serviceAccountName,omitempty"`
	ServiceAccount                *ServiceAccount `json:"serviceAccount,omitempty"`
	NodeName                      *NodeName       `json:"nodeName,omitempty"`
	SecurityContext               *EmptyDirClass  `json:"securityContext,omitempty"`
	SchedulerName                 *SchedulerName  `json:"schedulerName,omitempty"`
	Tolerations                   []Toleration    `json:"tolerations,omitempty"`
	Priority                      *int64          `json:"priority,omitempty"`
	EnableServiceLinks            *bool           `json:"enableServiceLinks,omitempty"`
	DNSConfig                     *DNSConfig      `json:"dnsConfig,omitempty"`
	HostAliases                   []HostAlias     `json:"hostAliases,omitempty"`
	Affinity                      *Affinity       `json:"affinity,omitempty"`
	NodeSelector                  *NodeSelector   `json:"nodeSelector,omitempty"`
	HostNetwork                   *bool           `json:"hostNetwork,omitempty"`
}

type Affinity struct {
	NodeAffinity *NodeAffinity `json:"nodeAffinity,omitempty"`
}

type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution *RequiredDuringSchedulingIgnoredDuringExecution `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

type RequiredDuringSchedulingIgnoredDuringExecution struct {
	NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms,omitempty"`
}

type NodeSelectorTerm struct {
	MatchFields []MatchField `json:"matchFields,omitempty"`
}

type MatchField struct {
	Key      *string    `json:"key,omitempty"`
	Operator *string    `json:"operator,omitempty"`
	Values   []NodeName `json:"values,omitempty"`
}

type Container struct {
	Name                     *string                   `json:"name,omitempty"`
	Image                    *string                   `json:"image,omitempty"`
	Ports                    []PodPort                 `json:"ports,omitempty"`
	Env                      []Env                     `json:"env,omitempty"`
	Resources                *Resources                `json:"resources,omitempty"`
	VolumeMounts             []VolumeMount             `json:"volumeMounts,omitempty"`
	TerminationMessagePath   *TerminationMessagePath   `json:"terminationMessagePath,omitempty"`
	TerminationMessagePolicy *TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty"`
	ImagePullPolicy          *ImagePullPolicy          `json:"imagePullPolicy,omitempty"`
	Args                     []string                  `json:"args,omitempty"`
	LivenessProbe            *NessProbe                `json:"livenessProbe,omitempty"`
	SecurityContext          *SecurityContext          `json:"securityContext,omitempty"`
	ReadinessProbe           *NessProbe                `json:"readinessProbe,omitempty"`
	Command                  []string                  `json:"command,omitempty"`
}

type Env struct {
	Name      *string    `json:"name,omitempty"`
	Value     *string    `json:"value,omitempty"`
	ValueFrom *ValueFrom `json:"valueFrom,omitempty"`
}

type ValueFrom struct {
	FieldRef *FieldRef `json:"fieldRef,omitempty"`
}

type FieldRef struct {
	APIVersion *string `json:"apiVersion,omitempty"`
	FieldPath  *string `json:"fieldPath,omitempty"`
}

type NessProbe struct {
	HTTPGet             *HTTPGet `json:"httpGet,omitempty"`
	InitialDelaySeconds *int64   `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      *int64   `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       *int64   `json:"periodSeconds,omitempty"`
	SuccessThreshold    *int64   `json:"successThreshold,omitempty"`
	FailureThreshold    *int64   `json:"failureThreshold,omitempty"`
}

type HTTPGet struct {
	Path   *string `json:"path,omitempty"`
	Port   *int64  `json:"port,omitempty"`
	Scheme *string `json:"scheme,omitempty"`
}

type PodPort struct {
	ContainerPort *int64    `json:"containerPort,omitempty"`
	Protocol      *Protocol `json:"protocol,omitempty"`
	Name          *string   `json:"name,omitempty"`
	HostPort      *int64    `json:"hostPort,omitempty"`
}

type Resources struct {
	Limits   *Limits `json:"limits,omitempty"`
	Requests *Limits `json:"requests,omitempty"`
}

type Limits struct {
	CPU    *string `json:"cpu,omitempty"`
	Memory *string `json:"memory,omitempty"`
}

type SecurityContext struct {
	Capabilities             *Capabilities `json:"capabilities,omitempty"`
	ReadOnlyRootFilesystem   *bool         `json:"readOnlyRootFilesystem,omitempty"`
	AllowPrivilegeEscalation *bool         `json:"allowPrivilegeEscalation,omitempty"`
	RunAsUser                *int64        `json:"runAsUser,omitempty"`
	RunAsNonRoot             *bool         `json:"runAsNonRoot,omitempty"`
	RunAsGroup               *int64        `json:"runAsGroup,omitempty"`
}

type Capabilities struct {
	Add  []string `json:"add,omitempty"`
	Drop []string `json:"drop,omitempty"`
}

type VolumeMount struct {
	Name      *SecretNameEnum `json:"name,omitempty"`
	MountPath *MountPath      `json:"mountPath,omitempty"`
	ReadOnly  *bool           `json:"readOnly,omitempty"`
}

type DNSConfig struct {
	Options []Option `json:"options,omitempty"`
}

type Option struct {
	Name  *OptionName `json:"name,omitempty"`
	Value *string     `json:"value,omitempty"`
}

type HostAlias struct {
	IP        *string  `json:"ip,omitempty"`
	Hostnames []string `json:"hostnames,omitempty"`
}

type InitContainer struct {
	Name                     *InitContainerName        `json:"name,omitempty"`
	Image                    *ImageStr                 `json:"image,omitempty"`
	Command                  []Command                 `json:"command,omitempty"`
	Args                     []string                  `json:"args,omitempty"`
	Resources                *Resources                `json:"resources,omitempty"`
	VolumeMounts             []VolumeMount             `json:"volumeMounts,omitempty"`
	TerminationMessagePath   *TerminationMessagePath   `json:"terminationMessagePath,omitempty"`
	TerminationMessagePolicy *TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty"`
	ImagePullPolicy          *ImagePullPolicy          `json:"imagePullPolicy,omitempty"`
}

type NodeSelector struct {
	KubernetesIoOS *string `json:"kubernetes.io/os,omitempty"`
}

type EmptyDirClass struct {
}

type Toleration struct {
	Key               *Key      `json:"key,omitempty"`
	Operator          *Operator `json:"operator,omitempty"`
	Effect            *Effect   `json:"effect,omitempty"`
	TolerationSeconds *int64    `json:"tolerationSeconds,omitempty"`
}

type Volume struct {
	Name                  *SecretNameEnum        `json:"name,omitempty"`
	EmptyDir              *EmptyDirClass         `json:"emptyDir,omitempty"`
	Secret                *Secret                `json:"secret,omitempty"`
	ConfigMap             *ConfigMap             `json:"configMap,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`
}

type ConfigMap struct {
	Name        *string         `json:"name,omitempty"`
	DefaultMode *int64          `json:"defaultMode,omitempty"`
	Items       []ConfigMapItem `json:"items,omitempty"`
}

type ConfigMapItem struct {
	Key  *string `json:"key,omitempty"`
	Path *string `json:"path,omitempty"`
}

type PersistentVolumeClaim struct {
	ClaimName *string `json:"claimName,omitempty"`
}

type Secret struct {
	SecretName  *SecretNameEnum `json:"secretName,omitempty"`
	DefaultMode *int64          `json:"defaultMode,omitempty"`
}

type PodStatusClass struct {
	Phase                 *Phase                `json:"phase,omitempty"`
	Conditions            []Condition           `json:"conditions,omitempty"`
	HostIP                *NodeName             `json:"hostIP,omitempty"`
	PodIP                 *string               `json:"podIP,omitempty"`
	PodIPS                []PodIP               `json:"podIPs,omitempty"`
	StartTime             *string               `json:"startTime,omitempty"`
	InitContainerStatuses []InitContainerStatus `json:"initContainerStatuses,omitempty"`
	ContainerStatuses     []ContainerStatus     `json:"containerStatuses,omitempty"`
	QosClass              *QosClass             `json:"qosClass,omitempty"`
	Message               *string               `json:"message,omitempty"`
	Reason                *StatusReason         `json:"reason,omitempty"`
}

type PodCondition struct {
	Type               *Type            `json:"type,omitempty"`
	Status             *StatusEnum      `json:"status,omitempty"`
	LastProbeTime      interface{}      `json:"lastProbeTime"`
	LastTransitionTime *string          `json:"lastTransitionTime,omitempty"`
	Reason             *ConditionReason `json:"reason,omitempty"`
	Message            *string          `json:"message,omitempty"`
}

type ContainerStatus struct {
	Name         *string         `json:"name,omitempty"`
	State        *PurpleState    `json:"state,omitempty"`
	LastState    *LastStateClass `json:"lastState,omitempty"`
	Ready        *bool           `json:"ready,omitempty"`
	RestartCount *int64          `json:"restartCount,omitempty"`
	Image        *string         `json:"image,omitempty"`
	ImageID      *string         `json:"imageID,omitempty"`
	ContainerID  *string         `json:"containerID,omitempty"`
	Started      *bool           `json:"started,omitempty"`
}

type LastStateClass struct {
	Terminated *Terminated `json:"terminated,omitempty"`
}

type Terminated struct {
	ExitCode    *int64            `json:"exitCode,omitempty"`
	Reason      *TerminatedReason `json:"reason,omitempty"`
	StartedAt   *string           `json:"startedAt,omitempty"`
	FinishedAt  *string           `json:"finishedAt,omitempty"`
	ContainerID *string           `json:"containerID,omitempty"`
}

type PurpleState struct {
	Running    *RunningClass `json:"running,omitempty"`
	Waiting    *Waiting      `json:"waiting,omitempty"`
	Terminated *Terminated   `json:"terminated,omitempty"`
}

type RunningClass struct {
	StartedAt *string `json:"startedAt,omitempty"`
}

type Waiting struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

type InitContainerStatus struct {
	Name         *InitContainerName `json:"name,omitempty"`
	State        *LastStateClass    `json:"state,omitempty"`
	LastState    *EmptyDirClass     `json:"lastState,omitempty"`
	Ready        *bool              `json:"ready,omitempty"`
	RestartCount *int64             `json:"restartCount,omitempty"`
	Image        *string            `json:"image,omitempty"`
	ImageID      *string            `json:"imageID,omitempty"`
	ContainerID  *string            `json:"containerID,omitempty"`
}

type PodIP struct {
	IP *string `json:"ip,omitempty"`
}

type PodInfoMetadata struct {
	SelfLink        *string `json:"selfLink,omitempty"`
	ResourceVersion *string `json:"resourceVersion,omitempty"`
}

type ServiceAccount string

const (
	Coredns                           ServiceAccount = "coredns"
	KubeStateMetrics                  ServiceAccount = "kube-state-metrics"
	MetricsServer                     ServiceAccount = "metrics-server"
	NginxIngressServiceaccount        ServiceAccount = "nginx-ingress-serviceaccount"
	ServiceAccountDefault             ServiceAccount = "default"
	ServiceAccountKubernetesDashboard ServiceAccount = "kubernetes-dashboard"
)

const (
	NamespaceDefault             Namespace = "default"
	NamespaceKubernetesDashboard Namespace = "kubernetes-dashboard"
)

type APIVersion string

const (
	AppsV1  APIVersion = "apps/v1"
	BatchV1 APIVersion = "batch/v1"
)

const (
	DaemonSet  Kind = "DaemonSet"
	Job        Kind = "Job"
	ReplicaSet Kind = "ReplicaSet"
)

type NodeName string

const (
	The19216810111 NodeName = "192.168.10.111"
	The1921681020  NodeName = "192.168.10.20"
	The1921681021  NodeName = "192.168.10.21"
	The1921681022  NodeName = "192.168.10.22"
	The1921681023  NodeName = "192.168.10.23"
	The1921681024  NodeName = "192.168.10.24"
	The1921681032  NodeName = "192.168.10.32"
	The1921681063  NodeName = "192.168.10.63"
)

type ImagePullPolicy string

const (
	IfNotPresent          ImagePullPolicy = "IfNotPresent"
	ImagePullPolicyAlways ImagePullPolicy = "Always"
)

type TerminationMessagePath string

const (
	DevTerminationLog TerminationMessagePath = "/dev/termination-log"
)

type TerminationMessagePolicy string

const (
	File TerminationMessagePolicy = "File"
)

type MountPath string

const (
	Certs                                   MountPath = "/certs"
	EtcCoredns                              MountPath = "/etc/coredns"
	EtcInfoRuihuaFusionStandard             MountPath = "/etc/info.ruihua.fusion.standard/"
	OptSrcMainResourcesStaticUploadImgs     MountPath = "/opt/src/main/resources/static/upload/imgs"
	SkywalkingAgent                         MountPath = "/skywalking/agent"
	Tmp                                     MountPath = "/tmp"
	USRShareNginxHTML                       MountPath = "/usr/share/nginx/html"
	USRSkywalkingAgent                      MountPath = "/usr/skywalking/agent"
	VarRunSecretsKubernetesIoServiceaccount MountPath = "/var/run/secrets/kubernetes.io/serviceaccount"
)

type SecretNameEnum string

const (
	Application                          SecretNameEnum = "application"
	ConfigVolume                         SecretNameEnum = "config-volume"
	CorednsTokenGmtj7                    SecretNameEnum = "coredns-token-gmtj7"
	DefaultToken5K774                    SecretNameEnum = "default-token-5k774"
	DefaultToken7S6F2                    SecretNameEnum = "default-token-7s6f2"
	Imgs                                 SecretNameEnum = "imgs"
	KubeStateMetricsToken7Pk6F           SecretNameEnum = "kube-state-metrics-token-7pk6f"
	KubernetesDashboardCerts             SecretNameEnum = "kubernetes-dashboard-certs"
	KubernetesDashboardTokenQb5R2        SecretNameEnum = "kubernetes-dashboard-token-qb5r2"
	MetricsServerTokenLsq96              SecretNameEnum = "metrics-server-token-lsq96"
	NameSkywalkingAgent                  SecretNameEnum = "skywalking-agent"
	NginxIngressServiceaccountTokenXgf9N SecretNameEnum = "nginx-ingress-serviceaccount-token-xgf9n"
	TmpDir                               SecretNameEnum = "tmp-dir"
	TmpVolume                            SecretNameEnum = "tmp-volume"
)

type OptionName string

const (
	Ndots OptionName = "ndots"
)

type DNSPolicy string

const (
	ClusterFirst DNSPolicy = "ClusterFirst"
	// Default      DNSPolicy = "Default"
)

type Command string

const (
	Sh Command = "sh"
)

type ImageStr string

const (
	ImageCloudLocal4000SidecarSkywalking890 ImageStr = "image.cloud.local:4000/sidecar-skywalking:8.9.0"
)

type InitContainerName string

const (
	SkywalkingAgentSidecar InitContainerName = "skywalking-agent-sidecar"
)

type RestartPolicy string

const (
	OnFailure           RestartPolicy = "OnFailure"
	RestartPolicyAlways RestartPolicy = "Always"
)

type SchedulerName string

const (
	DefaultScheduler SchedulerName = "default-scheduler"
)

type Effect string

const (
	NoExecute  Effect = "NoExecute"
	NoSchedule Effect = "NoSchedule"
)

type Key string

const (
	CriticalAddonsOnly             Key = "CriticalAddonsOnly"
	NodeKubernetesIoDiskPressure   Key = "node.kubernetes.io/disk-pressure"
	NodeKubernetesIoMemoryPressure Key = "node.kubernetes.io/memory-pressure"
	NodeKubernetesIoNotReady       Key = "node.kubernetes.io/not-ready"
	NodeKubernetesIoPIDPressure    Key = "node.kubernetes.io/pid-pressure"
	NodeKubernetesIoUnreachable    Key = "node.kubernetes.io/unreachable"
	NodeKubernetesIoUnschedulable  Key = "node.kubernetes.io/unschedulable"
	NodeRoleKubernetesIoMaster     Key = "node-role.kubernetes.io/master"
)

type Operator string

const (
	Exists Operator = "Exists"
)

type ConditionReason string

const (
	ContainersNotReady ConditionReason = "ContainersNotReady"
	PodCompleted       ConditionReason = "PodCompleted"
)

// type StatusEnum string

// const (
// 	False StatusEnum = "False"
// 	True  StatusEnum = "True"
// )

// type Type string

const (
	ContainersReady Type = "ContainersReady"
	Initialized     Type = "Initialized"
	PodScheduled    Type = "PodScheduled"
	// Ready           Type = "Ready"
)

type TerminatedReason string

const (
	Completed TerminatedReason = "Completed"
	Error     TerminatedReason = "Error"
)

type Phase string

const (
	Failed    Phase = "Failed"
	Running   Phase = "Running"
	Succeeded Phase = "Succeeded"
)

type QosClass string

const (
	BestEffort QosClass = "BestEffort"
	Burstable  QosClass = "Burstable"
	Guaranteed QosClass = "Guaranteed"
)

type StatusReason string

const (
	Evicted StatusReason = "Evicted"
)
