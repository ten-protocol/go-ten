import { siteMetadata } from "../../../tools/bridge-frontend/src/lib/siteMetadata";
import { DocumentInterface } from "../lib/types/common";

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
