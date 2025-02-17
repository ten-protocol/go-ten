import React from "react";
import Layout from "@/src/components/layouts/default-layout";
import EmptyState from "@repo/ui/components/common/empty-state";
import { Button } from "@repo/ui/components/shared/button";
import { useRouter } from "next/router";

const AddressDetails = () => {
  const { push } = useRouter();

  return (
    <Layout>
      <EmptyState
        title="Address Details"
        description="Coming soon..."
        imageSrc="/assets/images/clock.png"
        imageAlt="Clock"
        action={<Button onClick={() => push("/")}>Go Home</Button>}
      />
    </Layout>
  );
};

export default AddressDetails;

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
