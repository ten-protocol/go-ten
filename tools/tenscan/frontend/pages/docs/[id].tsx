import Layout from "@/src/components/layouts/default-layout";
import { useRouter } from "next/router";
import React from "react";
import DocumentComponent from "@repo/ui/components/common/document";
import { useDocumentService } from "@repo/ui/services/useGeneralService";
import EmptyState from "@repo/ui/components/common/empty-state";
import { Button } from "@repo/ui/components/shared/button";

const DocumentPage = () => {
  const { query, push } = useRouter();
  const { id } = query as { id: string };

  const { data: document, isLoading } = useDocumentService(id);

  return (
    <Layout>
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {document ? (
          <DocumentComponent document={document} isLoading={isLoading} />
        ) : (
          <EmptyState
            title="Document not found"
            description="The document you are looking for does not exist."
            action={<Button onClick={() => push("/")}>Go home</Button>}
          />
        )}
      </div>
    </Layout>
  );
};

export default DocumentPage;

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
