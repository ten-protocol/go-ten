export interface PaginationInterface {
    page: number;
    perPage: number;
    total: number;
    totalPages: number;
}

export interface ResponseDataInterface<T> {
    result: T;
    item: T;
    message: string;
    pagination?: PaginationInterface;
    success: boolean;
}

export type NavLink = {
    label: string;
    href?: string;
    isDropdown?: boolean;
    isExternal?: boolean;
    subNavLinks?: NavLink[];
};

export type Environment = 'uat-testnet' | 'sepolia-testnet' | 'dev-testnet' | 'local-testnet';
