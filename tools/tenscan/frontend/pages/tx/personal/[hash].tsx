import Layout from "@/src/components/layouts/default-layout";
import EmptyState from "@/src/components/modules/common/empty-state";
import { Button } from "@/src/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@/src/components/ui/card";
import { Skeleton } from "@/src/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";
import { fetchPersonalTxnByHash } from "@/api/transactions";
import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import { PersonalTxnDetailsComponent } from "@/src/components/modules/personal/personal-txn-details";

export default function TransactionDetails() {
  const router = useRouter();
  const { provider } = useWalletConnection();
  const { hash } = router.query;

  const { data: transactionDetails, isLoading } = useQuery({
    queryKey: ["personalTxnData", hash],
    queryFn: () => fetchPersonalTxnByHash(provider, hash as string),
    enabled: !!provider && !!hash,
  });

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
            <PersonalTxnDetailsComponent
              transactionDetails={transactionDetails}
            />
          </CardContent>
        </Card>
      ) : (
        <EmptyState
          title="Transaction not found"
          description="The transaction you are looking for does not exist."
          action={
            <Button onClick={() => router.push("/personal")}>Go back</Button>
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
