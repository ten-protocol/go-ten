import { fetchTransactionByHash } from "@/api/transactions";
import Layout from "@/src/components/layouts/default-layout";
import { TransactionDetailsComponent } from "@/src/components/modules/transactions/transaction-details";
import EmptyState from "@repo/ui/components/common/empty-state";
import { Button } from "@repo/ui/components/shared/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@repo/ui/components/shared/card";
import { Skeleton } from "@repo/ui/components/shared/skeleton";
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
      <Card className="col-span-3">
        {isLoading ? (
          <>
            <Skeleton className="h-10 w-100" />
            <Skeleton className="h-10 w-100" />
            <Skeleton className="h-10 w-100" />
            <Skeleton className="h-10 w-100" />
            <Skeleton className="h-10 w-100" />
            <Skeleton className="h-10 w-100" />
          </>
        ) : transactionDetails ? (
          <>
            <CardHeader>
              <CardTitle>Transaction Details</CardTitle>
            </CardHeader>
            <CardContent>
              <TransactionDetailsComponent
                transactionDetails={transactionDetails}
              />
            </CardContent>
          </>
        ) : (
          <EmptyState
            title="Transaction not found"
            description="The transaction you are looking for does not exist."
            action={
              <Button onClick={() => router.push("/transactions")}>
                Go back
              </Button>
            }
            className="p-8"
          />
        )}
      </Card>
    </Layout>
  );
}

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
