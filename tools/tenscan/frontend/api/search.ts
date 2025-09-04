import {SearchResponse} from "@/src/types/interfaces/SearchInterfaces";
import {httpRequest} from "@/api/index";
import {pathToUrl} from "@/src/routes/router";
import {apiRoutes} from "@/src/routes";
import {ResponseDataInterface} from "@repo/ui/lib/types/common";

export const searchRecords = async (
    query: string
): Promise<ResponseDataInterface<SearchResponse>> => {
    return await httpRequest<ResponseDataInterface<SearchResponse>>({
        method: "get",
        url: pathToUrl(apiRoutes.getSearchResults),
        searchParams: { query },
    });
};