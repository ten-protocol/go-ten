import { tenGatewayAddress } from "../lib/constants";
import axios, { AxiosInstance, AxiosRequestConfig } from "axios";

type HttpMethod = "get" | "post" | "put" | "patch" | "delete";

interface HttpOptions {
  method?: HttpMethod;
  url: string;
  data?: Record<string, any>;
  params?: Record<string, any>;
  headers?: Record<string, any>;
  timeout?: number;
  responseType?:
    | "json"
    | "arraybuffer"
    | "blob"
    | "document"
    | "text"
    | undefined;
  searchParams?: Record<string, any>;
}

const baseConfig: AxiosRequestConfig = {
  baseURL: tenGatewayAddress,
  timeout: 10000,
};

export const https: AxiosInstance = axios.create(baseConfig);

export const httpRequest = async <ResponseData>(
  options: HttpOptions,
  config: AxiosRequestConfig = {}
): Promise<ResponseData> => {
  const {
    method = "get",
    url,
    data,
    params,
    headers,
    timeout,
    responseType,
    searchParams,
  } = options;
  let query = "";
  if (searchParams) {
    const filteredParams = Object.fromEntries(
      Object.entries(searchParams).filter(
        ([, value]) => value !== undefined && value !== null && value !== ""
      )
    );
    if (Object.keys(filteredParams).length) {
      query = new URLSearchParams(filteredParams).toString();
    }
  }

  const httpConfig: AxiosRequestConfig = {
    method,
    url: query ? `${url}?${query}` : url,
    data,
    params,
    headers: { ...(headers || {}) },
    timeout,
    responseType: responseType,
    ...config,
  };
  const response = await https(httpConfig);
  return response.data as ResponseData;
};
