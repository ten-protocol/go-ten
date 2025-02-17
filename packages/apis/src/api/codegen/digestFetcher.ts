import { resolveFetchError, UnacceptedError, resolveUrl } from "../client";
import { FetchError, FetchOptions, ofetch } from "ofetch";
import { authSubject } from "../../libs/observables/auth";
import { destr } from "destr";
import { ApiError } from "../response-types";
import { DigestContext } from "./digestContext";

export type ErrorWrapper<TError> =
  | TError
  | { status: "unknown"; payload: string };

export type DigestFetcherOptions<TBody, THeaders, TQueryParams, TPathParams> = {
  url: string;
  method: string;
  body?: TBody;
  headers?: THeaders;
  queryParams?: TQueryParams;
  pathParams?: TPathParams;
  signal?: AbortSignal;
} & DigestContext["fetcherOptions"];

export async function digestFetch<
  TData,
  TError,
  TBody extends Record<string, unknown> | FormData | undefined | null,
  THeaders extends {},
  TQueryParams extends {},
  TPathParams extends {},
>(
  params: DigestFetcherOptions<TBody, THeaders, TQueryParams, TPathParams>,
): Promise<TData> {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...params.headers,
  };

  try {
    return UnacceptedError.interceptResponse(
      await baseFetcherFn({
        ...params,
        query: params.queryParams,
        headers,
        url: resolveUrl(params.url, params.queryParams, params.pathParams),
      }),
    );
  } catch (error) {
    return Promise.reject(resolveFetchError(error));
  }
}

export async function baseFetcherFn<TResponse>(
  params: { url: string } & FetchOptions,
): Promise<TResponse> {
  try {
    const { url, ...rest } = params;
    const request = api<TResponse>(url, rest as FetchOptions<"json">);

    // @ts-expect-error
    return Promise.race([
      new Promise((_, rej) =>
        setTimeout(
          () => rej(ApiError({ message: "Request Timed out" })),
          30000,
        ),
      ),
      request,
    ]);
  } catch (err) {
    const xhrErr = err as FetchError;
    if (xhrErr instanceof FetchError) {
      if (xhrErr?.response?.status === 401) {
        authSubject.next({ type: "unauthorized" });
      }
    }
    throw err;
  }
}

function resolveEnvVariable(
  record: Record<string, string>,
  env_variable: string,
): string {
  if (env_variable in record) {
    return String(record[env_variable]);
  }
  return "--";
}

const url = resolveEnvVariable(import.meta.env || process.env, "VITE_API_URL");

const api = ofetch.create({
  baseURL: url,
  parseResponse: destr,
});
