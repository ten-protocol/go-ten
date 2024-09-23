import { fetchTransactionByHash } from "@/api/transactions";
import Layout from "@/src/components/layouts/default-layout";
import { TransactionDetailsComponent } from "@/src/components/modules/transactions/transaction-details";
import EmptyState from "@repo/ui/common/empty-state";
import { Button } from "@repo/ui/shared/button";
import { Card, CardHeader, CardTitle, CardContent } from "@repo/ui/shared/card";
import { Skeleton } from "@repo/ui/shared/skeleton";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";

export default function TransactionDetails() {
  const router = useRouter();
  const { hash } = router.query;

  const { data, isLoading } = useQuery({
    queryKey: ["transactionDetails", hash],
    queryFn: () => fetchTransactionByHash(hash as string),
  });

  const transactionDetails = data?.item;

  return (
    <Layout>
      {isLoading ? (
        <Skeleton className="h-full w-full" />
      ) : transactionDetails ? (
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Transaction Details</CardTitle>
          </CardHeader>
          <CardContent>
            <TransactionDetailsComponent
              transactionDetails={transactionDetails}
            />
          </CardContent>
        </Card>
      ) : (
        <EmptyState
          title="Transaction not found"
          description="The transaction you are looking for does not exist."
          action={
            <Button onClick={() => router.push("/transactions")}>
              Go back
            </Button>
          }
        />
      )}
    </Layout>
  );
}

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
