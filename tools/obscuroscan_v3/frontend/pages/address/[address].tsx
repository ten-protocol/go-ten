import Layout from "@/src/components/layouts/default-layout";
import EmptyState from "@/src/components/modules/common/empty-state";
import { Button } from "@/src/components/ui/button";
import { useRouter } from "next/router";
import React from "react";

const AddressDetails = () => {
  //   const { query } = useRouter();
  //   const { address } = query;

  return (
    <Layout>
      <EmptyState
        title="Address Details"
        description="Coming soon..."
        image="/assets/images/clock.png"
        action={<Button>Go Home</Button>}
      />
    </Layout>
  );
};

export default AddressDetails;
