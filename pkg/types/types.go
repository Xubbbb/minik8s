package types

type VolumeMount struct {
	Name      string `yaml:"name" json:"name"`
	MountPath string `yaml:"mountPath" json:"mountPath"`
	ReadOnly  bool   `yaml:"readOnly" json:"readOnly"`
}

type ResourceList struct {
	CPU    int64 `yaml:"cpu" json:"cpu"`
	Memory int64 `yaml:"memory" json:"memory"`
}

type ResourceRequirements struct {

	// Limits describes the maximum amount of compute resources allowed.
	Limits ResourceList `yaml:"limits" json:"limits"`

	// Requests describes the minimum amount of compute resources required.
	Requests ResourceList `yaml:"requests" json:"requests"`
}

type EnvVar struct {
	// Required: Name of the environment variable.
	// When the RelaxedEnvironmentVariableValidation feature gate is disabled, this must consist of alphabetic characters,
	// digits, '_', '-', or '.', and must not start with a digit.
	// When the RelaxedEnvironmentVariableValidation feature gate is enabled,
	// this may contain any printable ASCII characters except '='.
	Name string
	// Optional: no more than one of the following may be specified.
	// Optional: Defaults to ""; variable references $(VAR_NAME) are expanded
	// using the previously defined environment variables in the container and
	// any service environment variables.  If a variable cannot be resolved,
	// the reference in the input string will be unchanged.  Double $$ are
	// reduced to a single $, which allows for escaping the $(VAR_NAME)
	// syntax: i.e. "$$(VAR_NAME)" will produce the string literal
	// "$(VAR_NAME)".  Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// +optional
	Value string
	// Optional: Specifies a source the value of this var should come from.
	// +optional
	// ValueFrom *EnvVarSource
}

type Protocol string

const (
	// ProtocolTCP is the TCP protocol.
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP is the UDP protocol.
	ProtocolUDP Protocol = "UDP"
	// ProtocolSCTP is the SCTP protocol.
	ProtocolSCTP Protocol = "SCTP"
)

// ContainerPort represents a network port in a single container
type ContainerPort struct {
	// Optional: If specified, this must be an IANA_SVC_NAME  Each named port
	// in a pod must have a unique name.
	// +optional
	Name string
	// Optional: If specified, this must be a valid port number, 0 < x < 65536.
	// If HostNetwork is specified, this must match ContainerPort.
	// +optional
	HostPort int32
	// Required: This must be a valid port number, 0 < x < 65536.
	ContainerPort int32
	// Required: Supports "TCP", "UDP" and "SCTP"
	// +optional
	Protocol Protocol
	// Optional: What host IP to bind the external port to.
	// +optional
	HostIP string
}

type Container struct {
	Name string `yaml:"name" json:"name"`

	Image string `yaml:"image" json:"image"`

	Command []string `yaml:"command" json:"command"`

	Args []string `yaml:"args" json:"args"`

	WorkingDir string `yaml:"workingDir" json:"workingDir"`

	Ports []ContainerPort `yaml:"ports" json:"ports"`

	Env []EnvVar `yaml:"env" json:"env"`

	Resources ResourceRequirements `yaml:"resources" json:"resources"`

	VolumeMounts []VolumeMount `yaml:"volumeMounts" json:"volumeMounts"`
}

type ContainerState string

const (
	// ContainerStateRunning means the container is currently running.
	ContainerStateRunning ContainerState = "Running"
	// ContainerStateTerminated means the container has exited with failure.
	ContainerStateTerminated ContainerState = "Terminated"
	// ContainerStateWaiting means the container is waiting to run.
	ContainerStateWaiting ContainerState = "Waiting"
)

type PodPhase string

// These are the valid statuses of pods.
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	// Deprecated: It isn't being set since 2015 (74da3b14b0c0f658b3bb8d2def5094686d0e9095)
	PodUnknown PodPhase = "Unknown"
)

type ContainerStatus struct {
	Name string `yaml:"name" json:"name"`

	State ContainerState `yaml:"state" json:"state"`

	Ready bool `yaml:"ready" json:"ready"`

	Started string `yaml:"started" json:"started"`

	Finished string `yaml:"finished" json:"finished"`
}

type PodStatus struct {
	Phase PodPhase `yaml:"phase" json:"phase"`

	PodIP string `yaml:"podIP" json:"podIP"`

	ContainerStatuses []ContainerStatus `yaml:"containerStatuses" json:"containerStatuses"`

	CPUUsage float64 `yaml:"cpuUsage" json:"cpuUsage"`

	MemoryUsage float64 `yaml:"memoryUsage" json:"memoryUsage"`
}
