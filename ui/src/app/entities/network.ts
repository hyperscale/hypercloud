import { IPAM } from './ipam';

export interface NetworkContainer {
    EndpointID?: string;
    IPv4Address?: string;
    IPv6Address?: string;
    MacAddress?: string;
    Name?: string;
}

export interface Network {
    Attachable?: boolean;
    Created?: string;
    ConfigOnly?: boolean;
    Options?: { [key: string]: string };
    Containers?: { [key: string]: NetworkContainer };
    Driver?: string;
    EnableIPv6?: boolean;
    IPAM?: IPAM;
    Id?: string;
    Ingress?: boolean;
    Internal?: boolean;
    Labels?: { [key: string]: string };
    Name?: string;
    Scope?: string;
}
