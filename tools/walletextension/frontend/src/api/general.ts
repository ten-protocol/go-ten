import { ResponseDataInterface } from '@/types/interfaces';
import { httpRequest } from './index';
import { pathToUrl } from '@/routes/router';
import { apiRoutes } from '@/routes';

export const fetchTestnetStatus = async (): Promise<ResponseDataInterface<any>> => {
    return await httpRequest<ResponseDataInterface<any>>({
        method: 'get',
        url: pathToUrl(apiRoutes.getHealthStatus),
    });
};
