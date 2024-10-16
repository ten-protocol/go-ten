import { Data, Match } from "effect";

interface ApiError {
  readonly _tag: "ApiError"; // the tag
  message: string;
  error?: unknown;
}

interface ValidationError {
  readonly _tag: "ValidationError"; // the tag
  messages: Record<string, string>;
}

interface UnknownError {
  readonly _tag: "UnknownError"; // the tag
  message: string;
  value: unknown;
}

export interface Okay<T> {
  readonly _tag: "Okay"; // the tag
  value: T;
}

export const ApiError = Data.tagged<ApiError>("ApiError");
export const ValidationError = Data.tagged<ValidationError>("ValidationError");
export const UnknownError = Data.tagged<UnknownError>("UnknownError");

export const Okay = <T>(data: T) => {
  const Constructor = Data.case<Okay<T>>();
  return Constructor({ _tag: "Okay", value: data });
};

export type ApiErrors = ValidationError | ApiError | UnknownError;

export type SafeRes<T> = ApiErrors | Okay<T>;

export const ResponseType = Match.type<ApiErrors | Okay<unknown>>();
