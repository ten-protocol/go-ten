import { safeStr } from "@repo/shared/src/data.helpers";
import { Match } from "effect";
import { FetchError } from "ofetch";
import { has, identity, is } from "ramda";
import {
  ApiError,
  ApiErrors,
  Okay,
  ResponseType,
  SafeRes,
  UnknownError,
} from "./response-types";

export const resolveUrl = (
  url: string,
  queryParams: Record<string, string> = {},
  pathParams: Record<string, string> = {},
) => {
  let query = new URLSearchParams(queryParams).toString();
  if (query) query = `?${query}`;

  return url.replace(/\{\w*\}/g, (key) => pathParams[key.slice(1, -1)]) + query;
};

export async function safeFetchResponse<T>(
  promise: Promise<T>,
): Promise<SafeRes<T>> {
  try {
    const value = await promise;
    return resolveFetchResponse<T>(value);
  } catch (err: unknown) {
    return resolveFetchError(err);
  }
}

export const resolveFetchResponse = <T>(response: T): SafeRes<T> => {
  const hasError = has("error", response) && response.error === true;
  const hasMessage = has("msg", response);

  const messageIsStr = hasMessage && typeof response?.msg === "string";

  if (UnacceptedError.validate(response)) {
    return UnacceptedError.respond(response);
  }

  if (Array.isArray(response)) return Okay(response);

  if (hasError && hasMessage && is(Object, response?.msg))
    // @ts-expect-error Still trying to response structure
    return ValidationError(response?.msg);

  if (hasError && hasMessage && messageIsStr)
    return ApiError({
      message: response.msg as string,

      error: response.error,
    });

  if (response === undefined || response === null) {
    return UnknownError({
      message: "Something went wrong",
      value: response,
    });
  }

  return Okay(response);
};

// biome-ignore lint/suspicious/noExplicitAny: Error type must be any
export function resolveFetchError(error: any): ApiErrors {
  if (
    ["Network Error", "NetworkError"].some((str) =>
      safeStr(error?.message).includes(str),
    )
  ) {
    return ApiError({
      message: "Seems like you're offline. Please check your network",
    });
  }

  return handleErrorByType(error);
}

const handleErrorByType = ResponseType.pipe(
  Match.tag("ApiError", identity),
  Match.tag("ValidationError", identity),
  Match.tag("UnknownError", identity),
  Match.orElse((err) => handleRandomError(err)),
);

function handleRandomError(error: unknown) {
  const err = error as FetchError;
  if (err instanceof FetchError) {
    const { status = 400 } = guessRequestError(err);
    const reason_for_failure = err?.data?.message || err.message;

    if (status === 404)
      return ApiError({ message: "404: Resource not found", error: err });

    if (401 === status) {
      return ApiError({
        message: "Unable to process request. You're unauthorized",
        error: err,
      });
    }

    if (403 === status) {
      return ApiError({
        message: "You do not have permission to perform this operation",
        error: err,
      });
    }

    if (status >= 400 && status < 500)
      return ApiError({ message: reason_for_failure, error: err });

    if (status >= 500)
      return ApiError({ message: reason_for_failure, error: err });
  }

  const err_msg = has("message", error) ? safeStr(error.message) : "";

  if (err_msg.includes("<no response>")) {
    return ApiError({
      message: "Something new wrong. Server didn't respond",
    });
  }

  return ApiError({ message: err_msg });

  function guessRequestError(err: FetchError) {
    return err instanceof FetchError
      ? err
      : { status: 0, statusText: "UNKNOWN" };
  }
}

export type ResponseResolver = {
  validate: (data: unknown) => boolean;
  respond: <T>(data: T) => ApiErrors;
  interceptResponse: <T>(data: T) => T;
};

export const UnacceptedError: ResponseResolver = {
  validate(response: unknown) {
    return (
      has("accepted", response) &&
      has("message", response) &&
      response.accepted === false
    );
  },

  respond(response: unknown) {
    return ApiError({
      // @ts-expect-error
      message: response?.message ?? "Unknown Error...",
      error: response,
    });
  },

  interceptResponse(response) {
    if (!this.validate(response)) return response;

    throw this.respond(response);
  },
};
