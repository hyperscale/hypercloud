export interface Version {
    Index?: number;
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
    State: string;
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
