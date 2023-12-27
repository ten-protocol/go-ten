import Layout from "@/src/components/layouts/default-layout";
import EmptyState from "@/src/components/modules/common/empty-state";
import { Button } from "@/src/components/ui/button";
import { useRouter } from "next/router";
import React from "react";

const TransactionDetails = () => {
  const { push } = useRouter();

  return (
    <Layout>
      <EmptyState
        title="Transaction Details"
        description="Coming soon..."
        imageSrc="/assets/images/clock.png"
        imageAlt="Clock"
        action={<Button onClick={() => push("/")}>Go Home</Button>}
      />
    </Layout>
  );
};

export default TransactionDetails;
