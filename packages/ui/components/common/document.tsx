import React from "react";
import Custom404Error from "./404";
import Spinner from "../shared/spinner";
import {
  DocumentContentInterface,
  DocumentInterface,
} from "../../lib/types/common";

type DocumentProps = {
  document: DocumentInterface;
  isLoading: boolean;
};

const DocumentComponent = ({ document, isLoading }: DocumentProps) => {
  if (isLoading) {
    return <Spinner />;
  }

  if (!document?.title) {
    return <Custom404Error customPageTitle="Document" />;
  }

  return (
    <>
      <div className="mb-8 text-center">
        <h1 className="text-4xl font-extrabold mb-6">{document.title}</h1>
        <p className="text-sm text-muted-foreground">{document.subHeading}</p>
      </div>
      <div className="prose prose-lg prose-primary">
        {document.content &&
          document.content.map(
            (section: DocumentContentInterface, index: number) => (
              <div key={index} className="mb-8">
                <h2 className="mb-2">{section.heading}</h2>
                {section.content &&
                  section.content.map((paragraph: string, index: number) => (
                    <div
                      key={index}
                      dangerouslySetInnerHTML={{ __html: paragraph }}
                    ></div>
                  ))}
              </div>
            )
          )}
      </div>
    </>
  );
};

export default DocumentComponent;
