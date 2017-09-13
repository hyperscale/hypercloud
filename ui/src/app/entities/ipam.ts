
export interface IPAMConfig {
    Gateway?: string;
    Subnet?: string;
}

export interface IPAM {
    Config?: IPAMConfig[];
    Driver?: string;
    Options?: any;
}
