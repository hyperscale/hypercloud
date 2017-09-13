import { Version } from './version';

export interface JoinTokens {
    Worker: string;
    Manager: string;
}

export interface SwarmOrchestration {
    TaskHistoryRetentionLimit: number;
}

export interface SwarmRaft {
    SnapshotInterval: number;
    LogEntriesForSlowFollowers: number;
    HeartbeatTick: number;
    ElectionTick: number;
}

export interface SwarmDispatcher {
    HeartbeatPeriod: number;
}

export interface SwarmCAConfig {
    NodeCertExpiry: number;
}

export interface SwarmTaskDefaults {
}

export interface SwarmSpec {
    Name: string;
    Orchestration: SwarmOrchestration;
    Raft: SwarmRaft;
    Dispatcher: SwarmDispatcher;
    CAConfig: SwarmCAConfig;
    TaskDefaults: SwarmTaskDefaults;
}

export interface Swarm {
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    JoinTokens: JoinTokens;
    Spec: SwarmSpec;
}
