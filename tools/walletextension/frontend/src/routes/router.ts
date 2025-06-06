import { compile } from 'path-to-regexp';

export const pathToUrl = (path: string, params: object = {}) => {
    return compile(path)(params);
};
