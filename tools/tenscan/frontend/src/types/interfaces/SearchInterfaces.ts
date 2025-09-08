export interface SearchResponse {
    ResultsData: SearchResult[];
    Total: number;
}

export interface SearchResult {
    type: string;                    // "rollup", "batch", "transaction"
    hash: string;
    height?: bigint;                 // batches only, optional
    sequence?: bigint;               // batches only, optional
    timestamp: number;               // rollup and batches
    extraData: Record<string, any>;  // contains all the batch/ rollup/ tx data if found
}