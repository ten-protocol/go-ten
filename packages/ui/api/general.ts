import { httpRequest } from ".";
import { siteMetadata } from "../../../tools/bridge-frontend/src/lib/siteMetadata";
import { DocumentInterface, ResponseDataInterface } from "../lib/types/common";
import { apiRoutes } from "../routes";
import { pathToUrl } from "../routes/router";

export const fetchTestnetStatus = async (): Promise<
  ResponseDataInterface<boolean>
> => {
  return await httpRequest<ResponseDataInterface<boolean>>({
    method: "get",
    url: pathToUrl(apiRoutes.getHealthStatus),
  });
};
export const fetchDocument = async (id: string): Promise<DocumentInterface> => {
  const response = await fetch(`/docs/${id}.json`);
  const data = await response.json();

  const processedData = {
    title: data.title,
    subHeading: data.subHeading,
    content: data.content.map((item: any) => {
      return {
        heading: item.heading,
        content: item.content.map((paragraph: any) => {
          return paragraph.replace(/siteMetadata.email/g, siteMetadata.email);
        }),
      };
    }),
  };

  return processedData;
};
