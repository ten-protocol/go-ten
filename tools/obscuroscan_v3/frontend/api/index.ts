import { apiHost } from "@/src/lib/constants";
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
  download?: boolean;
  searchParams?: Record<string, any>;
}

const baseConfig: AxiosRequestConfig = {
  baseURL: apiHost,
  timeout: 10000,
};

const https: AxiosInstance = axios.create(baseConfig);

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
  try {
    const response = await https(httpConfig);
    return response.data as ResponseData;
  } catch (error) {
    handleHttpError(error);
    throw error;
  }
};

// Centralized error handling function
const handleHttpError = (error: any) => {
  // if the error is a server error (status code 5xx) before handling
  if (isAxiosError(error) && error.response && error.response.status >= 500) {
    console.error("Server error:", error);
  } else {
    // other errors
    console.error("An error occurred:", error);
  }
};

// Type guard to check if the error is an AxiosError
const isAxiosError = (error: any): error is import("axios").AxiosError => {
  return error.isAxiosError === true;
};
