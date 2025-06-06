/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"k8s.io/kops/util/pkg/reflectutils"
)

// NetworkingSpec allows selection and configuration of a networking plugin
type NetworkingSpec struct {
	NetworkID              string              `json:"-"`
	NetworkCIDR            string              `json:"-"`
	AdditionalNetworkCIDRs []string            `json:"-"`
	Subnets                []ClusterSubnetSpec `json:"-"`
	TagSubnets             *bool               `json:"-"`
	Topology               *TopologySpec       `json:"-"`
	EgressProxy            *EgressProxySpec    `json:"-"`
	NonMasqueradeCIDR      string              `json:"-"`
	PodCIDR                string              `json:"-"`
	ServiceClusterIPRange  string              `json:"-"`
	IsolateControlPlane    *bool               `json:"-"`

	Classic    *ClassicNetworkingSpec    `json:"classic,omitempty"`
	Kubenet    *KubenetNetworkingSpec    `json:"kubenet,omitempty"`
	External   *ExternalNetworkingSpec   `json:"external,omitempty"`
	CNI        *CNINetworkingSpec        `json:"cni,omitempty"`
	Kopeio     *KopeioNetworkingSpec     `json:"kopeio,omitempty"`
	Weave      *WeaveNetworkingSpec      `json:"weave,omitempty"`
	Flannel    *FlannelNetworkingSpec    `json:"flannel,omitempty"`
	Calico     *CalicoNetworkingSpec     `json:"calico,omitempty"`
	Canal      *CanalNetworkingSpec      `json:"canal,omitempty"`
	KubeRouter *KuberouterNetworkingSpec `json:"kuberouter,omitempty"`
	Romana     *RomanaNetworkingSpec     `json:"romana,omitempty"`
	AmazonVPC  *AmazonVPCNetworkingSpec  `json:"amazonvpc,omitempty"`
	Cilium     *CiliumNetworkingSpec     `json:"cilium,omitempty"`
	LyftVPC    *LyftVPCNetworkingSpec    `json:"lyftvpc,omitempty"`
	GCP        *GCPNetworkingSpec        `json:"gce,omitempty"`
	Kindnet    *KindnetNetworkingSpec    `json:"kindnet,omitempty"`
}

func (s *NetworkingSpec) IsEmpty() bool {
	return s.ConfiguredOptions().Len() == 0
}

// ConfiguredOptions returns the set of networking options that are configured (non-nil)
// in the struct.  We only expect a single option to be configured.
func (s *NetworkingSpec) ConfiguredOptions() sets.Set[string] {
	options, err := reflectutils.FindSetFields(s, "classic", "kubenet", "external", "cni", "kopeio", "weave", "flannel", "calico", "canal", "kuberouter", "romana", "amazonvpc", "cilium", "lyftvpc", "gce", "kindnet")
	if err != nil {
		klog.Fatalf("error getting configured options: %v", err)
	}
	return options
}

// ClassicNetworkingSpec is the specification of classic networking mode, integrated into kubernetes.
// Support been removed since Kubernetes 1.4.
type ClassicNetworkingSpec struct{}

// KubenetNetworkingSpec is the specification for kubenet networking, largely integrated but intended to replace classic
type KubenetNetworkingSpec struct{}

// ExternalNetworkingSpec is the specification for networking that is implemented by a user-provided Daemonset that uses the Kubenet kubelet networking plugin.
type ExternalNetworkingSpec struct{}

// CNINetworkingSpec is the specification for networking that is implemented by a user-provided Daemonset, which uses the CNI kubelet networking plugin.
type CNINetworkingSpec struct {
	UsesSecondaryIP bool `json:"usesSecondaryIP,omitempty"`
}

// KopeioNetworkingSpec declares that we want Kopeio networking
type KopeioNetworkingSpec struct{}

// WeaveNetworkingSpec declares that we want Weave networking
type WeaveNetworkingSpec struct {
	MTU         *int32 `json:"mtu,omitempty"`
	ConnLimit   *int32 `json:"connLimit,omitempty"`
	NoMasqLocal *int32 `json:"noMasqLocal,omitempty"`

	// MemoryRequest memory request of weave container. Default 200Mi
	MemoryRequest *resource.Quantity `json:"memoryRequest,omitempty"`
	// CPURequest CPU request of weave container. Default 50m
	CPURequest *resource.Quantity `json:"cpuRequest,omitempty"`
	// MemoryLimit memory limit of weave container. Default 200Mi
	MemoryLimit *resource.Quantity `json:"memoryLimit,omitempty"`
	// CPULimit CPU limit of weave container.
	CPULimit *resource.Quantity `json:"cpuLimit,omitempty"`
	// NetExtraArgs are extra arguments that are passed to weave-kube.
	NetExtraArgs string `json:"netExtraArgs,omitempty"`

	// NPCMemoryRequest memory request of weave npc container. Default 200Mi
	NPCMemoryRequest *resource.Quantity `json:"npcMemoryRequest,omitempty"`
	// NPCCPURequest CPU request of weave npc container. Default 50m
	NPCCPURequest *resource.Quantity `json:"npcCPURequest,omitempty"`
	// NPCMemoryLimit memory limit of weave npc container. Default 200Mi
	NPCMemoryLimit *resource.Quantity `json:"npcMemoryLimit,omitempty"`
	// NPCCPULimit CPU limit of weave npc container
	NPCCPULimit *resource.Quantity `json:"npcCPULimit,omitempty"`
	// NPCExtraArgs are extra arguments that are passed to weave-npc.
	NPCExtraArgs string `json:"npcExtraArgs,omitempty"`

	// Version specifies the Weave container image tag. The default depends on the kOps version.
	Version string `json:"version,omitempty"`
}

// FlannelNetworkingSpec declares that we want Flannel networking
type FlannelNetworkingSpec struct {
	// Backend is the backend overlay type we want to use (vxlan or udp)
	Backend string `json:"backend,omitempty"`
	// DisableTxChecksumOffloading is unused.
	// +k8s:conversion-gen=false
	DisableTxChecksumOffloading bool `json:"disableTxChecksumOffloading,omitempty"`
	// IptablesResyncSeconds sets resync period for iptables rules, in seconds
	IptablesResyncSeconds *int32 `json:"iptablesResyncSeconds,omitempty"`
}

// CalicoNetworkingSpec declares that we want Calico networking
type CalicoNetworkingSpec struct {
	// Registry overrides the Calico container image registry.
	Registry string `json:"registry,omitempty"`
	// Version overrides the Calico container image tag.
	Version string `json:"version,omitempty"`

	// AllowIPForwarding enable ip_forwarding setting within the container namespace.
	// (default: false)
	AllowIPForwarding bool `json:"allowIPForwarding,omitempty"`
	// AWSSrcDstCheck enables/disables ENI source/destination checks (AWS IPv4 only)
	// Options: Disable (default for IPv4), Enable, or DoNothing
	AWSSrcDstCheck string `json:"awsSrcDstCheck,omitempty"`
	// BPFEnabled enables the eBPF dataplane mode.
	BPFEnabled bool `json:"bpfEnabled,omitempty"`
	// BPFExternalServiceMode controls how traffic from outside the cluster to NodePorts and ClusterIPs is handled.
	// In Tunnel mode, packet is tunneled from the ingress host to the host with the backing pod and back again.
	// In DSR mode, traffic is tunneled to the host with the backing pod and then returned directly;
	// this requires a network that allows direct return.
	// Default: Tunnel (other options: DSR)
	BPFExternalServiceMode string `json:"bpfExternalServiceMode,omitempty"`
	// BPFKubeProxyIptablesCleanupEnabled controls whether Felix will clean up the iptables rules
	// created by the Kubernetes kube-proxy; should only be enabled if kube-proxy is not running.
	BPFKubeProxyIptablesCleanupEnabled bool `json:"bpfKubeProxyIptablesCleanupEnabled,omitempty"`
	// BPFLogLevel controls the log level used by the BPF programs. The logs are emitted
	// to the BPF trace pipe, accessible with the command tc exec BPF debug.
	// Default: Off (other options: Info, Debug)
	BPFLogLevel string `json:"bpfLogLevel,omitempty"`
	// ChainInsertMode controls whether Felix inserts rules to the top of iptables chains, or
	// appends to the bottom. Leaving the default option is safest to prevent accidentally
	// breaking connectivity. Default: 'insert' (other options: 'append')
	ChainInsertMode string `json:"chainInsertMode,omitempty"`
	// CPURequest CPU request of Calico container. Default: 100m
	CPURequest *resource.Quantity `json:"cpuRequest,omitempty"`
	// CrossSubnet is deprecated as of kOps 1.22 and has no effect
	CrossSubnet *bool `json:"crossSubnet,omitempty"`
	// EncapsulationMode specifies the network packet encapsulation protocol for Calico to use,
	// employing such encapsulation at the necessary scope per the related CrossSubnet field. In
	// "ipip" mode, Calico will use IP-in-IP encapsulation as needed. In "vxlan" mode, Calico will
	// encapsulate packets as needed using the VXLAN scheme.
	// Options: ipip (default) or vxlan
	EncapsulationMode string `json:"encapsulationMode,omitempty"`
	// IPIPMode determines when to use IP-in-IP encapsulation for the default Calico IPv4 pool.
	// It is conveyed to the "calico-node" daemon container via the CALICO_IPV4POOL_IPIP
	// environment variable. EncapsulationMode must be set to "ipip".
	// Options: "CrossSubnet", "Always", or "Never".
	// Default: "CrossSubnet" if EncapsulationMode is "ipip", "Never" otherwise.
	IPIPMode string `json:"ipipMode,omitempty"`
	// IPv4AutoDetectionMethod configures how Calico chooses the IP address used to route
	// between nodes.  This should be set when the host has multiple interfaces
	// and it is important to select the interface used.
	// Options: "first-found" (default), "can-reach=DESTINATION",
	// "interface=INTERFACE-REGEX", or "skip-interface=INTERFACE-REGEX"
	IPv4AutoDetectionMethod string `json:"ipv4AutoDetectionMethod,omitempty"`
	// IPv6AutoDetectionMethod configures how Calico chooses the IP address used to route
	// between nodes.  This should be set when the host has multiple interfaces
	// and it is important to select the interface used.
	// Options: "first-found" (default), "can-reach=DESTINATION",
	// "interface=INTERFACE-REGEX", or "skip-interface=INTERFACE-REGEX"
	IPv6AutoDetectionMethod string `json:"ipv6AutoDetectionMethod,omitempty"`
	// IptablesBackend controls which variant of iptables binary Felix uses
	// Default: Auto (other options: Legacy, NFT)
	IptablesBackend string `json:"iptablesBackend,omitempty"`
	// LogSeverityScreen lets us set the desired log level. (Default: info)
	LogSeverityScreen string `json:"logSeverityScreen,omitempty"`
	// MTU to be set in the cni-network-config for calico.
	MTU *int32 `json:"mtu,omitempty"`
	// PrometheusMetricsEnabled can be set to enable the experimental Prometheus
	// metrics server (default: false)
	PrometheusMetricsEnabled bool `json:"prometheusMetricsEnabled,omitempty"`
	// PrometheusMetricsPort is the TCP port that the experimental Prometheus
	// metrics server should bind to (default: 9091)
	PrometheusMetricsPort int32 `json:"prometheusMetricsPort,omitempty"`
	// PrometheusGoMetricsEnabled enables Prometheus Go runtime metrics collection
	PrometheusGoMetricsEnabled bool `json:"prometheusGoMetricsEnabled,omitempty"`
	// PrometheusProcessMetricsEnabled enables Prometheus process metrics collection
	PrometheusProcessMetricsEnabled bool `json:"prometheusProcessMetricsEnabled,omitempty"`
	// MajorVersion is unused.
	// +k8s:conversion-gen=false
	MajorVersion string `json:"majorVersion,omitempty"`
	// TyphaPrometheusMetricsEnabled enables Prometheus metrics collection from Typha
	// (default: false)
	TyphaPrometheusMetricsEnabled bool `json:"typhaPrometheusMetricsEnabled,omitempty"`
	// TyphaPrometheusMetricsPort is the TCP port the typha Prometheus metrics server
	// should bind to (default: 9093)
	TyphaPrometheusMetricsPort int32 `json:"typhaPrometheusMetricsPort,omitempty"`
	// TyphaReplicas is the number of replicas of Typha to deploy
	TyphaReplicas int32 `json:"typhaReplicas,omitempty"`
	// VXLANMode determines when to use VXLAN encapsulation for the default Calico IPv4 pool.
	// It is conveyed to the "calico-node" daemon container via the CALICO_IPV4POOL_VXLAN
	// environment variable. EncapsulationMode must be set to "vxlan".
	// Options: "CrossSubnet", "Always", or "Never".
	// Default: "CrossSubnet" if EncapsulationMode is "vxlan", "Never" otherwise.
	VXLANMode string `json:"vxlanMode,omitempty"`
	// WireguardEnabled enables WireGuard encryption for all on-the-wire pod-to-pod traffic
	// (default: false)
	WireguardEnabled bool `json:"wireguardEnabled,omitempty"`
}

// CanalNetworkingSpec declares that we want Canal networking
type CanalNetworkingSpec struct {
	// ChainInsertMode controls whether Felix inserts rules to the top of iptables chains, or
	// appends to the bottom. Leaving the default option is safest to prevent accidentally
	// breaking connectivity. Default: 'insert' (other options: 'append')
	ChainInsertMode string `json:"chainInsertMode,omitempty"`
	// CPURequest CPU request of Canal container. Default: 100m
	CPURequest *resource.Quantity `json:"cpuRequest,omitempty"`
	// DefaultEndpointToHostAction allows users to configure the default behaviour
	// for traffic between pod to host after calico rules have been processed.
	// Default: ACCEPT (other options: DROP, RETURN)
	DefaultEndpointToHostAction string `json:"defaultEndpointToHostAction,omitempty"`
	// DisableFlannelForwardRules configures Flannel to NOT add the
	// default ACCEPT traffic rules to the iptables FORWARD chain
	FlanneldIptablesForwardRules *bool `json:"disableFlannelForwardRules,omitempty"`
	// DisableTxChecksumOffloading is unused.
	// +k8s:conversion-gen=false
	DisableTxChecksumOffloading bool `json:"disableTxChecksumOffloading,omitempty"`
	// IptablesBackend controls which variant of iptables binary Felix uses
	// Default: Auto (other options: Legacy, NFT)
	IptablesBackend string `json:"iptablesBackend,omitempty"`
	// LogSeveritySys the severity to set for logs which are sent to syslog
	// Default: INFO (other options: DEBUG, WARNING, ERROR, CRITICAL, NONE)
	LogSeveritySys string `json:"logSeveritySys,omitempty"`
	// MTU to be set in the cni-network-config (default: 1500)
	MTU *int32 `json:"mtu,omitempty"`
	// PrometheusGoMetricsEnabled enables Prometheus Go runtime metrics collection
	PrometheusGoMetricsEnabled bool `json:"prometheusGoMetricsEnabled,omitempty"`
	// PrometheusMetricsEnabled can be set to enable the experimental Prometheus
	// metrics server (default: false)
	PrometheusMetricsEnabled bool `json:"prometheusMetricsEnabled,omitempty"`
	// PrometheusMetricsPort is the TCP port that the experimental Prometheus
	// metrics server should bind to (default: 9091)
	PrometheusMetricsPort int32 `json:"prometheusMetricsPort,omitempty"`
	// PrometheusProcessMetricsEnabled enables Prometheus process metrics collection
	PrometheusProcessMetricsEnabled bool `json:"prometheusProcessMetricsEnabled,omitempty"`
	// TyphaPrometheusMetricsEnabled enables Prometheus metrics collection from Typha
	// (default: false)
	TyphaPrometheusMetricsEnabled bool `json:"typhaPrometheusMetricsEnabled,omitempty"`
	// TyphaPrometheusMetricsPort is the TCP port the typha Prometheus metrics server
	// should bind to (default: 9093)
	TyphaPrometheusMetricsPort int32 `json:"typhaPrometheusMetricsPort,omitempty"`
	// TyphaReplicas is the number of replicas of Typha to deploy
	TyphaReplicas int32 `json:"typhaReplicas,omitempty"`
}

// KuberouterNetworkingSpec declares that we want Kube-router networking
type KuberouterNetworkingSpec struct{}

// RomanaNetworkingSpec declares that we want Romana networking
// Romana is deprecated as of kOps 1.18 and removed as of kOps 1.19.
type RomanaNetworkingSpec struct {
	// DaemonServiceIP is the Kubernetes Service IP for the romana-daemon pod
	DaemonServiceIP string `json:"daemonServiceIP,omitempty"`
	// EtcdServiceIP is the Kubernetes Service IP for the etcd backend used by Romana
	EtcdServiceIP string `json:"etcdServiceIP,omitempty"`
}

// AmazonVPCNetworkingSpec declares that we want Amazon VPC CNI networking
type AmazonVPCNetworkingSpec struct {
	// ImageName is the container image name to use.
	Image string `json:"imageName,omitempty"`
	// InitImageName is the init container image name to use.
	InitImage string `json:"initImageName,omitempty"`
	// NetworkPolicyAgentImage is the container image to use for the network policy agent
	NetworkPolicyAgentImage string `json:"networkPolicyAgentImage,omitempty"`
	// Env is a list of environment variables to set in the container.
	Env []EnvVar `json:"env,omitempty"`
}

const CiliumIpamEni = "eni"

type CiliumEncryptionType string

const (
	CiliumEncryptionTypeIPSec     CiliumEncryptionType = "ipsec"
	CiliumEncryptionTypeWireguard CiliumEncryptionType = "wireguard"
)

// CiliumNetworkingSpec declares that we want Cilium networking
type CiliumNetworkingSpec struct {
	// Registry overrides the default Cilium container registry (quay.io)
	Registry string `json:"registry,omitempty"`

	// Version is the version of the Cilium agent and the Cilium Operator.
	Version string `json:"version,omitempty"`

	// MemoryRequest memory request of Cilium agent + operator container. (default: 128Mi)
	MemoryRequest *resource.Quantity `json:"memoryRequest,omitempty"`
	// CPURequest CPU request of Cilium agent + operator container. (default: 25m)
	CPURequest *resource.Quantity `json:"cpuRequest,omitempty"`

	// AccessLog is unused.
	// +k8s:conversion-gen=false
	AccessLog string `json:"accessLog,omitempty"`
	// AgentLabels is unused.
	// +k8s:conversion-gen=false
	AgentLabels []string `json:"agentLabels,omitempty"`

	// AgentPrometheusPort is the port to listen to for Prometheus metrics.
	// Defaults to 9090.
	AgentPrometheusPort int `json:"agentPrometheusPort,omitempty"`
	// Metrics is a list of metrics to add or remove from the default list of metrics the agent exposes.
	Metrics []string `json:"metrics,omitempty"`

	// AllowLocalhost is unused.
	// +k8s:conversion-gen=false
	AllowLocalhost string `json:"allowLocalhost,omitempty"`
	// AutoIpv6NodeRoutes is unused.
	// +k8s:conversion-gen=false
	AutoIpv6NodeRoutes bool `json:"autoIpv6NodeRoutes,omitempty"`
	// BPFRoot is unused.
	// +k8s:conversion-gen=false
	BPFRoot string `json:"bpfRoot,omitempty"`
	// ChainingMode allows using Cilium in combination with other CNI plugins.
	// With Cilium CNI chaining, the base network connectivity and IP address management is managed
	// by the non-Cilium CNI plugin, but Cilium attaches eBPF programs to the network devices created
	// by the non-Cilium plugin to provide L3/L4 network visibility, policy enforcement and other advanced features.
	// Default: none
	ChainingMode string `json:"chainingMode,omitempty"`
	// ContainerRuntime is unused.
	// +k8s:conversion-gen=false
	ContainerRuntime []string `json:"containerRuntime,omitempty"`
	// ContainerRuntimeEndpoint is unused.
	// +k8s:conversion-gen=false
	ContainerRuntimeEndpoint map[string]string `json:"containerRuntimeEndpoint,omitempty"`
	// Debug runs Cilium in debug mode.
	Debug bool `json:"debug,omitempty"`
	// DebugVerbose is unused.
	// +k8s:conversion-gen=false
	DebugVerbose []string `json:"debugVerbose,omitempty"`
	// Device is unused.
	// +k8s:conversion-gen=false
	Device string `json:"device,omitempty"`
	// DisableConntrack is unused.
	// +k8s:conversion-gen=false
	DisableConntrack bool `json:"disableConntrack,omitempty"`
	// DisableEndpointCRD disables usage of CiliumEndpoint CRD.
	// Default: false
	DisableEndpointCRD bool `json:"disableEndpointCRD,omitempty"`
	// DisableIpv4 is unused.
	// +k8s:conversion-gen=false
	DisableIpv4 bool `json:"disableIpv4,omitempty"`
	// DisableK8sServices is unused.
	// +k8s:conversion-gen=false
	DisableK8sServices bool `json:"disableK8sServices,omitempty"`
	// EnablePolicy specifies the policy enforcement mode.
	// "default": Follows Kubernetes policy enforcement.
	// "always": Cilium restricts all traffic if no policy is in place.
	// "never": Cilium allows all traffic regardless of policies in place.
	// If unspecified, "default" policy mode will be used.
	EnablePolicy string `json:"enablePolicy,omitempty"`
	// EnableL7Proxy enables L7 proxy for L7 policy enforcement.
	// Default: true
	EnableL7Proxy *bool `json:"enableL7Proxy,omitempty"`
	// EnableLocalRedirectPolicy that enables pod traffic destined to an IP address and port/protocol
	// tuple or Kubernetes service to be redirected locally to backend pod(s) within a node, using eBPF.
	// https://docs.cilium.io/en/stable/network/kubernetes/local-redirect-policy/
	// Default: false
	EnableLocalRedirectPolicy *bool `json:"enableLocalRedirectPolicy,omitempty"`
	// EnableBPFMasquerade enables masquerading packets from endpoints leaving the host with BPF instead of iptables.
	// Default: false
	EnableBPFMasquerade *bool `json:"enableBPFMasquerade,omitempty"`
	// EnableEndpointHealthChecking enables connectivity health checking between virtual endpoints.
	// Default: true
	EnableEndpointHealthChecking *bool `json:"enableEndpointHealthChecking,omitempty"`
	// EnableTracing is unused.
	// +k8s:conversion-gen=false
	EnableTracing bool `json:"enableTracing,omitempty"`
	// EnablePrometheusMetrics enables the Cilium "/metrics" endpoint for both the agent and the operator.
	EnablePrometheusMetrics bool `json:"enablePrometheusMetrics,omitempty"`
	// EnableEncryption enables Cilium Encryption.
	// Default: false
	EnableEncryption bool `json:"enableEncryption,omitempty"`
	// EncryptionType specifies Cilium Encryption method ("ipsec", "wireguard").
	// Default: ipsec
	EncryptionType CiliumEncryptionType `json:"encryptionType,omitempty"`
	// NodeEncryption enables encryption for pure node to node traffic.
	// Default: false
	NodeEncryption bool `json:"nodeEncryption,omitempty"`
	// EnvoyLog is unused.
	// +k8s:conversion-gen=false
	EnvoyLog string `json:"envoyLog,omitempty"`
	// IdentityAllocationMode specifies in which backend identities are stored ("crd", "kvstore").
	// Default: crd
	IdentityAllocationMode string `json:"identityAllocationMode,omitempty"`
	// IdentityChangeGracePeriod specifies the duration to wait before using a changed identity.
	// Default: 5s
	IdentityChangeGracePeriod string `json:"identityChangeGracePeriod,omitempty"`
	// Ipv4ClusterCIDRMaskSize is unused.
	// +k8s:conversion-gen=false
	Ipv4ClusterCIDRMaskSize int `json:"ipv4ClusterCidrMaskSize,omitempty"`
	// Ipv4Node is unused.
	// +k8s:conversion-gen=false
	Ipv4Node string `json:"ipv4Node,omitempty"`
	// Ipv4Range is unused.
	// +k8s:conversion-gen=false
	Ipv4Range string `json:"ipv4Range,omitempty"`
	// Ipv4ServiceRange is unused.
	// +k8s:conversion-gen=false
	Ipv4ServiceRange string `json:"ipv4ServiceRange,omitempty"`
	// Ipv6ClusterAllocCidr is unused.
	// +k8s:conversion-gen=false
	Ipv6ClusterAllocCidr string `json:"ipv6ClusterAllocCidr,omitempty"`
	// Ipv6Node is unused.
	// +k8s:conversion-gen=false
	Ipv6Node string `json:"ipv6Node,omitempty"`
	// Ipv6Range is unused.
	// +k8s:conversion-gen=false
	Ipv6Range string `json:"ipv6Range,omitempty"`
	// Ipv6ServiceRange is unused.
	// +k8s:conversion-gen=false
	Ipv6ServiceRange string `json:"ipv6ServiceRange,omitempty"`
	// K8sAPIServer is unused.
	// +k8s:conversion-gen=false
	K8sAPIServer string `json:"k8sApiServer,omitempty"`
	// K8sKubeconfigPath is unused.
	// +k8s:conversion-gen=false
	K8sKubeconfigPath string `json:"k8sKubeconfigPath,omitempty"`
	// KeepBPFTemplates is unused.
	// +k8s:conversion-gen=false
	KeepBPFTemplates bool `json:"keepBpfTemplates,omitempty"`
	// KeepConfig is unused.
	// +k8s:conversion-gen=false
	KeepConfig bool `json:"keepConfig,omitempty"`
	// LabelPrefixFile is unused.
	// +k8s:conversion-gen=false
	LabelPrefixFile string `json:"labelPrefixFile,omitempty"`
	// Labels is unused.
	// +k8s:conversion-gen=false
	Labels []string `json:"labels,omitempty"`
	// LB is unused.
	// +k8s:conversion-gen=false
	LB string `json:"lb,omitempty"`
	// LibDir is unused.
	// +k8s:conversion-gen=false
	LibDir string `json:"libDir,omitempty"`
	// LogDrivers is unused.
	// +k8s:conversion-gen=false
	LogDrivers []string `json:"logDriver,omitempty"`
	// LogOpt is unused.
	// +k8s:conversion-gen=false
	LogOpt map[string]string `json:"logOpt,omitempty"`
	// Logstash is unused.
	// +k8s:conversion-gen=false
	Logstash bool `json:"logstash,omitempty"`
	// LogstashAgent is unused.
	// +k8s:conversion-gen=false
	LogstashAgent string `json:"logstashAgent,omitempty"`
	// LogstashProbeTimer is unused.
	// +k8s:conversion-gen=false
	LogstashProbeTimer uint32 `json:"logstashProbeTimer,omitempty"`
	// DisableMasquerade disables masquerading traffic to external destinations behind the node IP.
	Masquerade *bool `json:"disableMasquerade,omitempty"`
	// Nat46Range is unused.
	// +k8s:conversion-gen=false
	Nat46Range string `json:"nat46Range,omitempty"`
	// AgentPodAnnotations makes possible to add additional annotations to the cilium agent.
	// Default: none
	AgentPodAnnotations map[string]string `json:"agentPodAnnotations,omitempty"`
	// OperatorPodAnnotations makes possible to add additional annotations to cilium operator.
	// Default: none
	OperatorPodAnnotations map[string]string `json:"operatorPodAnnotations,omitempty"`
	// Pprof is unused.
	// +k8s:conversion-gen=false
	Pprof bool `json:"pprof,omitempty"`
	// PrefilterDevice is unused.
	// +k8s:conversion-gen=false
	PrefilterDevice string `json:"prefilterDevice,omitempty"`
	// PrometheusServeAddr is unused.
	// +k8s:conversion-gen=false
	PrometheusServeAddr string `json:"prometheusServeAddr,omitempty"`
	// Restore is unused.
	// +k8s:conversion-gen=false
	Restore bool `json:"restore,omitempty"`
	// SingleClusterRoute is unused.
	// +k8s:conversion-gen=false
	SingleClusterRoute bool `json:"singleClusterRoute,omitempty"`
	// SocketPath is unused.
	// +k8s:conversion-gen=false
	SocketPath string `json:"socketPath,omitempty"`
	// StateDir is unused.
	// +k8s:conversion-gen=false
	StateDir string `json:"stateDir,omitempty"`
	// TracePayloadLen is unused.
	// +k8s:conversion-gen=false
	TracePayloadLen int `json:"tracePayloadlen,omitempty"`
	// Tunnel specifies the Cilium tunnelling mode. Possible values are "vxlan", "geneve", or "disabled".
	// Default: vxlan
	Tunnel string `json:"tunnel,omitempty"`
	// EnableIpv6 is unused.
	// +k8s:conversion-gen=false
	EnableIpv6 bool `json:"enableipv6,omitempty"`
	// EnableIpv4 is unused.
	// +k8s:conversion-gen=false
	EnableIpv4 bool `json:"enableipv4,omitempty"`
	// MonitorAggregation sets the level of packet monitoring. Possible values are "low", "medium", or "maximum".
	// Default: medium
	MonitorAggregation string `json:"monitorAggregation,omitempty"`
	// BPFCTGlobalTCPMax is the maximum number of entries in the TCP CT table.
	// Default: 524288
	BPFCTGlobalTCPMax int `json:"bpfCTGlobalTCPMax,omitempty"`
	// BPFCTGlobalAnyMax is the maximum number of entries in the non-TCP CT table.
	// Default: 262144
	BPFCTGlobalAnyMax int `json:"bpfCTGlobalAnyMax,omitempty"`
	// BPFLBAlgorithm is the load balancing algorithm ("random", "maglev").
	// Default: random
	BPFLBAlgorithm string `json:"bpfLBAlgorithm,omitempty"`
	// BPFLBMaglevTableSize is the per service backend table size when going with Maglev (parameter M).
	// Default: 16381
	BPFLBMaglevTableSize string `json:"bpfLBMaglevTableSize,omitempty"`
	// BPFNATGlobalMax is the the maximum number of entries in the BPF NAT table.
	// Default: 524288
	BPFNATGlobalMax int `json:"bpfNATGlobalMax,omitempty"`
	// BPFNeighGlobalMax is the the maximum number of entries in the BPF Neighbor table.
	// Default: 524288
	BPFNeighGlobalMax int `json:"bpfNeighGlobalMax,omitempty"`
	// BPFPolicyMapMax is the maximum number of entries in endpoint policy map.
	// Default: 16384
	BPFPolicyMapMax int `json:"bpfPolicyMapMax,omitempty"`
	// BPFLBMapMax is the maximum number of entries in bpf lb service, backend and affinity maps.
	// Default: 65536
	BPFLBMapMax int `json:"bpfLBMapMax,omitempty"`
	// BPFLBSockHostNSOnly enables skipping socket LB for services when inside a pod namespace,
	// in favor of service LB at the pod interface. Socket LB is still used when in the host namespace.
	// Required by service mesh (e.g., Istio, Linkerd).
	// Default: false
	BPFLBSockHostNSOnly bool `json:"bpfLBSockHostNSOnly,omitempty"`
	// PreallocateBPFMaps reduces the per-packet latency at the expense of up-front memory allocation.
	// Default: true
	PreallocateBPFMaps bool `json:"preallocateBPFMaps,omitempty"`
	// SidecarIstioProxyImage is the regular expression matching compatible Istio sidecar istio-proxy
	// container image names.
	// Default: cilium/istio_proxy
	SidecarIstioProxyImage string `json:"sidecarIstioProxyImage,omitempty"`
	// ClusterName is the name of the cluster. It is only relevant when building a mesh of clusters.
	ClusterName string `json:"clusterName,omitempty"`
	// ClusterID is the ID of the cluster. It is only relevant when building a mesh of clusters.
	// Must be a number between 1 and 255.
	ClusterID uint8 `json:"clusterID,omitempty"`
	// ToFQDNsDNSRejectResponseCode sets the DNS response code for rejecting DNS requests.
	// Possible values are "nameError" or "refused".
	// Default: refused
	ToFQDNsDNSRejectResponseCode string `json:"toFqdnsDnsRejectResponseCode,omitempty"`
	// ToFQDNsEnablePoller replaces the DNS proxy-based implementation of FQDN policies
	// with the less powerful legacy implementation.
	// Default: false
	ToFQDNsEnablePoller bool `json:"toFqdnsEnablePoller,omitempty"`
	// ContainerRuntimeLabels is unused.
	// +k8s:conversion-gen=false
	ContainerRuntimeLabels string `json:"containerRuntimeLabels,omitempty"`
	// IPAM specifies the IP address allocation mode to use.
	// Possible values are "crd" and "eni".
	// "eni" will use AWS native networking for pods. Eni requires masquerade to be set to false.
	// "crd" will use CRDs for controlling IP address management.
	// "hostscope" will use hostscope IPAM mode.
	// "kubernetes" will use addersing based on node pod CIDR.
	// Default: "kubernetes".
	IPAM string `json:"ipam,omitempty"`
	// IPTablesRulesNoinstall disables installing the base IPTables rules used for masquerading and kube-proxy.
	// Default: false
	InstallIptablesRules *bool `json:"IPTablesRulesNoinstall,omitempty"`
	// AutoDirectNodeRoutes adds automatic L2 routing between nodes.
	// Default: false
	AutoDirectNodeRoutes bool `json:"autoDirectNodeRoutes,omitempty"`
	// EnableHostReachableServices configures Cilium to enable services to be
	// reached from the host namespace in addition to pod namespaces.
	// https://docs.cilium.io/en/v1.9/gettingstarted/host-services/
	// Default: false
	EnableHostReachableServices bool `json:"enableHostReachableServices,omitempty"`
	// EnableNodePort replaces kube-proxy with Cilium's BPF implementation.
	// Requires spec.kubeProxy.enabled be set to false.
	// Default: false
	EnableNodePort bool `json:"enableNodePort,omitempty"`
	// EtcdManagd installs an additional etcd cluster that is used for Cilium state change.
	// The cluster is operated by cilium-etcd-operator.
	// Default: false
	EtcdManaged bool `json:"etcdManaged,omitempty"`
	// EnableRemoteNodeIdentity enables the remote-node-identity.
	// Default: true
	EnableRemoteNodeIdentity *bool `json:"enableRemoteNodeIdentity,omitempty"`
	// EnableUnreachableRoutes enables unreachable routes on pod deletion.
	// Default: false
	EnableUnreachableRoutes *bool `json:"enableUnreachableRoutes,omitempty"`
	// CniExclusive configures whether to remove other CNI configuration files.
	// Default: true
	CniExclusive *bool `json:"cniExclusive,omitempty"`
	// Hubble configures the Hubble service on the Cilium agent.
	Hubble *HubbleSpec `json:"hubble,omitempty"`

	// RemoveCbrBridge is unused.
	// +k8s:conversion-gen=false
	RemoveCbrBridge bool `json:"removeCbrBridge,omitempty"`
	// RestartPods is unused.
	// +k8s:conversion-gen=false
	RestartPods bool `json:"restartPods,omitempty"`
	// ReconfigureKubelet is unused.
	// +k8s:conversion-gen=false
	ReconfigureKubelet bool `json:"reconfigureKubelet,omitempty"`
	// NodeInitBootstrapFile is unused.
	// +k8s:conversion-gen=false
	NodeInitBootstrapFile string `json:"nodeInitBootstrapFile,omitempty"`
	// CniBinPath is unused.
	// +k8s:conversion-gen=false
	CniBinPath string `json:"cniBinPath,omitempty"`
	// DisableCNPStatusUpdates determines if CNP NodeStatus updates will be sent to the Kubernetes api-server.
	DisableCNPStatusUpdates *bool `json:"disableCNPStatusUpdates,omitempty"`

	// EnableServiceTopology determine if cilium should use topology aware hints.
	EnableServiceTopology bool `json:"enableServiceTopology,omitempty"`

	// Ingress specifies the configuration for Cilium Ingress settings.
	Ingress *CiliumIngressSpec `json:"ingress,omitempty"`

	// GatewayAPI specifies the configuration for Cilium Gateway API settings.
	GatewayAPI *CiliumGatewayAPISpec `json:"gatewayAPI,omitempty"`
}

// CiliumIngressSpec configures Cilium Ingress settings.
type CiliumIngressSpec struct {
	// Enabled specifies whether Cilium Ingress is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// EnforceHttps specifies whether HTTPS enforcement is enabled for Ingress traffic.
	// Default: true
	EnforceHttps *bool `json:"enforceHttps,omitempty"`

	// EnableSecretsSync specifies whether synchronization of secrets is enabled.
	// Default: true
	EnableSecretsSync *bool `json:"enableSecretsSync,omitempty"`

	// LoadBalancerAnnotationPrefixes specifies annotation prefixes for Load Balancer configuration.
	// Default: "service.beta.kubernetes.io service.kubernetes.io cloud.google.com"
	LoadBalancerAnnotationPrefixes string `json:"loadBalancerAnnotationPrefixes,omitempty"`

	// DefaultLoadBalancerMode specifies the default load balancer mode.
	// Possible values: 'shared' or 'dedicated'
	// Default: dedicated
	DefaultLoadBalancerMode string `json:"defaultLoadBalancerMode,omitempty"`

	// SharedLoadBalancerServiceName specifies the name of the shared load balancer service.
	// Default: cilium-ingress
	SharedLoadBalancerServiceName string `json:"sharedLoadBalancerServiceName,omitempty"`
}

// CiliumGatewayAPISpec configures Cilium Gateway API settings.
type CiliumGatewayAPISpec struct {
	// Enabled specifies whether Cilium Gateway API is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// EnableSecretsSync specifies whether synchronization of secrets is enabled.
	// Default: true
	EnableSecretsSync *bool `json:"enableSecretsSync,omitempty"`
}

// HubbleSpec configures the Hubble service on the Cilium agent.
type HubbleSpec struct {
	// Enabled decides if Hubble is enabled on the agent or not
	Enabled *bool `json:"enabled,omitempty"`

	// Metrics is a list of metrics to collect. If empty or null, metrics are disabled.
	// See https://docs.cilium.io/en/stable/observability/metrics/#hubble-exported-metrics
	Metrics []string `json:"metrics,omitempty"`
}

// LyftVPCNetworkingSpec declares that we want to use the cni-ipvlan-vpc-k8s CNI networking.
// Lyft VPC is deprecated as of kOps 1.22 and removed as of kOps 1.23.
type LyftVPCNetworkingSpec struct {
	SubnetTags map[string]string `json:"subnetTags,omitempty"`
}

// GCPNetworkingSpec is the specification of GCP's native networking mode, using IP aliases.
type GCPNetworkingSpec struct{}

// KindnetNetworkingSpec configures Kindnet settings.
type KindnetNetworkingSpec struct {
	Version                      string                 `json:"version,omitempty"`
	NetworkPolicies              *bool                  `json:"networkPolicies,omitempty"`
	AdminNetworkPolicies         *bool                  `json:"adminNetworkPolicies,omitempty"`
	BaselineAdminNetworkPolicies *bool                  `json:"baselineAdminNetworkPolicies,omitempty"`
	DNSCaching                   *bool                  `json:"dnsCaching,omitempty"`
	NAT64                        *bool                  `json:"nat64,omitempty"`
	FastPathThreshold            *int32                 `json:"fastPathThreshold,omitempty"`
	Masquerade                   *KindnetMasqueradeSpec `json:"masquerade,omitempty"`
	LogLevel                     *int32                 `json:"logLevel,omitempty"`
}

// KindnetMasqueradeSpec configures Kindnet masquerading settings.
type KindnetMasqueradeSpec struct {
	Enabled            *bool    `json:"enabled,omitempty"`
	NonMasqueradeCIDRs []string `json:"nonMasqueradeCIDRs,omitempty"`
}
