export interface IRole {
    name: string;
    endpoint: {
        id: string,
        url: string
        method: string
    }
}