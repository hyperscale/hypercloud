import { Version } from './version';
import { ContainerSpec } from './container-spec';

export interface TaskResources {
    Limits: any;
    Reservations: any;
}

export interface TaskRestartPolicy {
    Condition: string;
    MaxAttempts: number;
}

export interface TaskPlacement {
    Constraints: string[];
}

export interface TaskSpec {
    ContainerSpec: ContainerSpec;
    Resources: TaskResources;
    RestartPolicy: TaskRestartPolicy;
    Placement: TaskPlacement;
}

export interface TaskContainerStatus {
    ContainerID: string;
}

export interface TaskStatus {
    Timestamp: Date;
    State: string;
    Message: string;
    ContainerStatus: TaskContainerStatus;
}

export interface TaskDriverConfiguration {
    Name: string;
}

export interface TaskDriver {
    Name: string;
}

export interface TaskIPAMOptionsConfigs {
    Subnet: string;
    Gateway: string;
}

export interface TaskIPAMOptions {
    Driver: TaskDriver;
    Configs?: TaskIPAMOptionsConfigs[];
}

export interface TaskNetworkSpec {
    Name: string;
    DriverConfiguration: TaskDriverConfiguration;
    IPAMOptions: TaskIPAMOptions;
}

export interface TaskDriverState {
    Name: string;
    Options: { [key:string]:string; };
}

export interface TaskNetwork {
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    Spec: TaskNetworkSpec;
    DriverState: TaskDriverState;
    IPAMOptions: TaskIPAMOptions;
}

export interface TaskNetworkAttachment {
    Network: TaskNetwork;
    Addresses: string[];
}

export interface Task {
    Name?: string;
    ContainerID?: string;
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    Spec: TaskSpec;
    ServiceID: string;
    Slot: number;
    NodeID: string;
    Status: TaskStatus;
    DesiredState: string;
    NetworksAttachments: TaskNetworkAttachment[];
}
