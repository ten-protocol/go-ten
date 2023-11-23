import Layout from "@/src/components/layouts/default-layout";
import Spinner from "@/src/components/ui/spinner";
import { useToast } from "@/src/components/ui/use-toast";
import { useRouter } from "next/router";
import React from "react";

type Document = {
  title: string;
  subHeading: string;
  content: {
    heading: string;
    content: string[];
  }[];
};

const Document = () => {
  const { toast } = useToast();
  const { query } = useRouter();
  const { id } = query;

  const [document, setDocument] = React.useState<Document>({} as Document);
  const [loading, setLoading] = React.useState<boolean>(false);

  const getDocument = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/docs/${id}.json`);
      const data = await response.json();
      setDocument(data);
    } catch (error) {
      toast({
        variant: "destructive",
        description: "Error fetching document",
      });
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    if (id) {
      getDocument();
    }
  }, [id]);

  return (
    <Layout>
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {!loading ? (
          <>
            <div className="mb-8 text-center">
              <h1 className="text-4xl font-extrabold mb-6">{document.title}</h1>
              <p className="text-sm text-muted-foreground">
                {document.subHeading}
              </p>
            </div>
            <div className="prose prose-lg prose-primary">
              {document.content &&
                document.content.map((section, index) => (
                  <div key={index} className="mb-8">
                    <h2 className="mb-2">{section.heading}</h2>
                    {section.content &&
                      section.content.map((paragraph, index) => (
                        <div
                          key={index}
                          dangerouslySetInnerHTML={{ __html: paragraph }}
                        ></div>
                      ))}
                  </div>
                ))}
            </div>
          </>
        ) : (
          <Spinner />
        )}
      </div>
    </Layout>
  );
};

export default Document;
