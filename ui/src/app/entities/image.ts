export interface Image {
    Id: string;
    ParentId?: string;
    RepoTags: string[];
    RepoDigests?: string[];
    Created: number;
    Size: number;
    VirtualSize: number;
    Labels?: { [key:string]:string; };
    Repository?: string;
    Tag?: string;
}
