export interface Version {
    Index?: number;
}

export interface Labels {
    [key: string]: string;
}

export interface NodeSpec {
    Role: string;
    Availability: string;
    Labels: { [key: string]: string };
}

export interface NodePlatform {
    Architecture: string;
    OS: string;
}

export interface NodeResources {
    NanoCPUs: number;
    MemoryBytes: number;
}

export interface NodePlugin {
    Type: string;
    Name: string;
}

export interface NodeEngine {
    EngineVersion: string;
    Plugins: NodePlugin[];
}

export interface NodeStatus {
    State: 'unknown' | 'down' | 'ready' | 'disconnected';
}

export interface NodeManagerStatus {
    Leader: boolean;
    Reachability: string;
    Addr: string;
}

export interface NodeDescription {
    Hostname: string;
    Platform: NodePlatform;
    Resources: NodeResources;
    Engine: NodeEngine;
    Status: NodeStatus;
    ManagerStatus: NodeManagerStatus;
}

export interface Node {
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    Spec: NodeSpec;
    Description: NodeDescription;
    Status?: NodeStatus;
}

export interface JoinTokens {
    Worker: string;
    Manager: string;
}

export interface Swarm {
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    JoinTokens: JoinTokens;
}

export interface InfoSwarmNode {
    NodeID: string;
    Addr: string;
}

export interface InfoSwarm {
    NodeID: string;
    NodeAddr: string;
    LocalNodeState: string;
    ControlAvailable: string;
    Error: string;
    RemoteManagers: InfoSwarmNode[];
    Nodes: number;
    Managers: number;
}

export interface Info {
    ID: string;
    Images: number;
    KernelVersion: string;
    OperatingSystem: string;
    OSType: string;
    Architecture: string;
    NCPU: number;
    MemTotal: number;
    ServerVersion: string;
    ClusterAdvertise: string;
    Swarm: InfoSwarm;
}

export interface Service {
    ID: string;
    Version: Version;
    CreatedAt: string;
    UpdatedAt: string;
    Spec: ServiceSpec;
    PreviousSpec: ServiceSpec;
    Endpoint: Endpoint;
    UpdateStatus?: UpdateStatus;
}

export interface UpdateStatus {
    State: string;
    StartedAt: string;
    CompletedAt: string;
    Message: string;
}

export interface ServiceSpec {
    Name: string;
    Labels: Labels;
    TaskTemplate: TaskSpec;
    Mode: Mode;
    UpdateConfig: UpdateConfig;
    EndpointSpec: EndpointSpecOrSpec;
}

export interface TaskSpec {
    ContainerSpec: ContainerSpec;
    Resources: ResourceRequirements;
    RestartPolicy: RestartPolicy;
    Placement: Placement;
    Networks?: NetworkAttachmentConfig[] | null;
    ForceUpdate: number;
    Runtime: string;
}

export interface ContainerSpec {
    Image: string;
    Labels: Labels;
    Privileges: Privileges;
    Isolation: string;
    Env: string[];
}


export interface Privileges {
    CredentialSpec?: null;
    SELinuxContext?: null;
}

export interface Resources {
    NanoCPUs: number;
    MemoryBytes: number;
    GenericResources: GenericResource[];
}

export interface GenericResource {
    NamedResourceSpec?: NamedGenericResource;
    DiscreteResourceSpec?: DiscreteGenericResource;
}

export interface DiscreteGenericResource {
    Kind: string;
    Value: string;
}

export interface NamedGenericResource {
    Kind: string;
    Value: string;
}

export interface ResourceRequirements {
    Limits?: Resources;
    Reservations?: Resources;
}

export interface RestartPolicy {
    Condition: string;
    MaxAttempts: number;
}

export interface Placement {
    Platforms?: (PlatformsEntity)[] | null;
}

export interface PlatformsEntity {
    Architecture?: string | null;
    OS: string;
}

export interface NetworkAttachmentConfig {
    Target: string;
    Aliases?: string[] | null;
    DriverOpts?: { [key: string]: string};
}

export interface Mode {
    Replicated?: Replicated;
    Global?: any;
}

export interface Replicated {
    Replicas: number;
}

export interface UpdateConfig {
    Parallelism: number;
    Delay?: number | null;
    FailureAction: string;
    MaxFailureRatio: number;
    Order: string;
}

export interface EndpointSpecOrSpec {
    Mode: string;
    Ports?: (PortConfig)[] | null;
}

export interface PortConfig {
    Name?: string;
    Protocol?: string;
    TargetPort?: number;
    PublishedPort?: number;
    PublishMode?: string;
}

export interface Endpoint {
    Spec: EndpointSpecOrSpec;
    Ports?: (PortConfig)[] | null;
    VirtualIPs?: (VirtualIPsEntity)[] | null;
}

export interface VirtualIPsEntity {
    NetworkID: string;
    Addr: string;
}


// ConainerStats

// ThrottlingData stores CPU throttling stats of one running container.
// Not used on Windows.
export interface ThrottlingData {
    // Number of periods with throttling active
    periods: number;
    // Number of periods when the container hits its throttling limit.
    throttled_periods: number;
    // Aggregate time the container was throttled for in nanoseconds.
    throttled_time: number;
}

// CPUUsage stores All CPU stats aggregated since container inception.
export interface CPUUsage {
    // Total CPU time consumed.
    // Units: nanoseconds (Linux)
    // Units: 100's of nanoseconds (Windows)
    total_usage: number;

    // Total CPU time consumed per core (Linux). Not used on Windows.
    // Units: nanoseconds.
    percpu_usage?: number[];

    // Time spent by tasks of the cgroup in kernel mode (Linux).
    // Time spent by all container processes in kernel mode (Windows).
    // Units: nanoseconds (Linux).
    // Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers.
    usage_in_kernelmode: number;

    // Time spent by tasks of the cgroup in user mode (Linux).
    // Time spent by all container processes in user mode (Windows).
    // Units: nanoseconds (Linux).
    // Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers
    usage_in_usermode: number;
}

// CPUStats aggregates and wraps all CPU related info of container
export interface CPUStats {
    // CPU Usage. Linux and Windows.
    cpu_usage: CPUUsage;

    // System Usage. Linux only.
    system_cpu_usage?: number;

    // Online CPUs. Linux only.
    online_cpus?: number;

    // Throttling Data. Linux only.
    throttling_data?: ThrottlingData;
}

// MemoryStats aggregates all memory stats since container inception on Linux.
// Windows returns stats for commit and private working set only.
export interface MemoryStats {
    // Linux Memory Stats

    // current res_counter usage for memory
    usage?: number;
    // maximum usage ever recorded.
    max_usage?: number;
    // TODO(vishh): Export these as stronger types.
    // all the stats exported via memory.stat.
    stats?: { [key: string]: number};
    // number of times memory usage hits limits.
    failcnt: number;
    limit: number;

    // Windows Memory Stats
    // See https://technet.microsoft.com/en-us/magazine/ff382715.aspx

    // committed bytes
    commitbytes?: number;
    // peak committed bytes
    commitpeakbytes?: number;
    // private working set
    privateworkingset?: number;
}

// BlkioStatEntry is one small entity to store a piece of Blkio stats
// Not used on Windows.
export interface BlkioStatEntry {
    major: number;
    minor: number;
    op: string;
    value: number;
}

// BlkioStats stores All IO service stats for data read and write.
// This is a Linux specific structure as the differences between expressing
// block I/O on Windows and Linux are sufficiently significant to make
// little sense attempting to morph into a combined structure.
export interface BlkioStats {
    // number of bytes transferred to and from the block device
    io_service_bytes_recursive: BlkioStatEntry[];
    io_serviced_recursive: BlkioStatEntry[];
    io_queue_recursive: BlkioStatEntry[];
    io_service_time_recursive: BlkioStatEntry[];
    io_wait_time_recursive: BlkioStatEntry[];
    io_merged_recursive: BlkioStatEntry[];
    io_time_recursive: BlkioStatEntry[];
    sectors_recursive: BlkioStatEntry[];
}

// StorageStats is the disk I/O stats for read/write on Windows.
export interface StorageStats {
    read_count_normalized?: number;
    read_size_bytes?: number;
    write_count_normalized?: number;
    write_size_bytes?: number;
}

// NetworkStats aggregates the network stats of one container
export interface NetworkStats {
    // Bytes received. Windows and Linux.
    rx_bytes: number;
    // Packets received. Windows and Linux.
    rx_packets: number;
    // Received errors. Not used on Windows. Note that we dont `omitempty` this
    // field as it is expected in the >=v1.21 API stats structure.
    rx_errors: number;
    // Incoming packets dropped. Windows and Linux.
    rx_dropped: number;
    // Bytes sent. Windows and Linux.
    tx_bytes: number;
    // Packets sent. Windows and Linux.
    tx_packets: number;
    // Sent errors. Not used on Windows. Note that we dont `omitempty` this
    // field as it is expected in the >=v1.21 API stats structure.
    tx_errors: number;
    // Outgoing packets dropped. Windows and Linux.
    tx_dropped: number;
    // Endpoint ID. Not used on Linux.
    endpoint_id?: string;
    // Instance ID. Not used on Linux.
    instance_id?: string;
}

// PidsStats contains the stats of a container's pids
export interface PidsStats {
    // Current is the number of pids in the cgroup
    current?: number;
    // Limit is the hard limit on the number of pids in the cgroup.
    // A "Limit" of 0 means that there is no limit.
    limit?: number;
}

// Stats is Ultimate aggregating all types of stats of one container
export interface Stats {
    // Common stats
    read: string;
    preread: string;

    // Linux specific stats, not populated on Windows.
    pids_stats?: PidsStats;
    blkio_stats?: BlkioStats;

    // Windows specific stats, not populated on Linux.
    num_procs: number;
    storage_stats?: StorageStats;

    // Shared stats
    cpu_stats?: CPUStats;

    precpu_stats?: CPUStats; // "Pre"="Previous"
    memory_stats?: MemoryStats;
}

// StatsJSON is newly used Networks
export interface StatsJSON extends Stats {
    name?: string;
    id?: string;

    // Networks request version >=1.21
    networks?: { [key: string]: NetworkStats };
}

// Actor describes something that generates events,
// like a container, or a network, or a volume.
// It has a defined name and a set or attributes.
// The container attributes are its labels, other actors
// can generate these attributes from other properties.
export interface Actor {
    ID: string;
    Attributes: { [key: string]: string };
}

// Message represents the information an event contains
export interface Message {
    // Deprecated information from JSONMessage.
    // With data only in container events.
    status?: string;
    id?: string;
    from?: string;

    Type: 'container' | 'daemon' | 'image' | 'network' | 'plugin' | 'volume' | 'service' | 'node' | 'secret' | 'config';
    Action: string;
    Actor: Actor;
    // Engine events are local scope. Cluster events are swarm scope.
    scope?: string;

    time?: number;
    timeNano?: number;
}

export interface Stack {
    Name: string;
    Services: number;
}

export interface VersionResponse {
    ApiVersion: string;
    Arch: string;
    BuildTime: string;
    Experimental: boolean;
    GitCommit: string;
    GoVersion: string;
    KernelVersion: string;
    MinAPIVersion: string;
    Os: string;
    Version: string;
}

export type TaskState =
    'new' |
    'allocated' |
    'pending' |
    'assigned' |
    'accepted' |
    'preparing' |
    'ready' |
    'starting' |
    'running' |
    'complete' |
    'shutdown' |
    'failed' |
    'rejected';

export interface ContainerStatus {
    ContainerID?: string;
    PID?: number;
    ExitCode?: number;
}

export interface PortStatus {
    Ports?: PortConfig[];
}

export interface TaskStatus {
    Timestamp?: string;
    State?: TaskState;
    Message?: string;
    Err?: string;
    ContainerStatus?: ContainerStatus;
    PortStatus?: PortStatus;
}

export interface Driver {
    Name?: string;
    Options: { [key: string]: string };
}


export interface IPAMConfig {
    Subnet?: string;
    Range?: string;
    Gateway?: string;
}

export interface IPAMOptions {
    Driver?: Driver;
    Configs?: IPAMConfig[];
}

export interface ConfigReference {
    Network: string;
}

export interface NetworkSpec {
    Name: string;
    Labels: { [key: string]: string };
    DriverConfiguration?: Driver;
    IPv6Enabled?: boolean;
    Internal?: boolean;
    Attachable?: boolean;
    Ingress?: boolean;
    IPAMOptions?: IPAMOptions;
    ConfigFrom?: ConfigReference;
    Scope?: string;
}

export interface Network {
    ID: string;
    Version: Version;
    CreatedAt: string;
    UpdatedAt: string;
    Spec?: NetworkSpec;
    DriverState?: Driver;
    IPAMOptions?: IPAMOptions;
}

export interface NetworkAttachment {
    Network?: Network;
    Addresses?: string[];
}

export interface Task {
    ID: string;
    Version: Version;
    CreatedAt: string;
    UpdatedAt: string;
    Name: string;
    Labels: { [key: string]: string };
    Spec?: TaskSpec;
    ServiceID?: string;
    Slot?: number;
    NodeID?: string;
    Status?: TaskStatus;
    DesiredState?: TaskState;
    NetworksAttachments?: NetworkAttachment[];
    GenericResources?: GenericResource[];
}

export interface ReplicasInfo {
    Running: number;
    Desired: number;
}

export interface ServiceInfo {
    Mode: string;
    Replicas: ReplicasInfo;
}
