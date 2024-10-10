import { is } from "ramda";

const EmptyPrimitives = Object.freeze({
  Array: [],
  Object: {
    __proto_: {
      type: "EmptyObject",
    },
  },
});

export const safeNum = (a: unknown, fallback = 0): number => {
  const value = Number(a);

  return !Object.is(NaN, value) ? value : fallback;
};

export const safeArray = <T>(a?: Array<T>): Array<T | never> =>
  Array.isArray(a) ? a : EmptyPrimitives.Array;

export const safeStr = (a: unknown, fallback = ""): string =>
  typeof a === "string" ? a : fallback;

const EmptyObject: Record<string, never> = Object.freeze({});

export const safeObj = <T>(
  obj: T,
): T extends Record<string, unknown> ? T : typeof EmptyObject => {
  // @ts-expect-error;
  return is(Object, obj) ? obj : EmptyObject;
};

export function safeInt(num: unknown, fallback = 0): number {
  const value = parseInt(num as string);

  return !Object.is(NaN, value) ? value : fallback;
}
