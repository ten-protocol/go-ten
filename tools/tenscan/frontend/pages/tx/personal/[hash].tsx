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
import ConnectWalletButton from "@/src/components/modules/common/connect-wallet";
import { ethereum } from "@/src/lib/utils";

export default function TransactionDetails() {
  const router = useRouter();
  const { provider, walletConnected } = useWalletConnection();
  const { hash } = router.query;

  const { data: transactionDetails, isLoading } = useQuery({
    queryKey: ["personalTxnData", hash],
    queryFn: () => fetchPersonalTxnByHash(provider, hash as string),
    enabled: !!provider && !!hash,
  });

  return (
    <Layout>
      {walletConnected ? (
        isLoading ? (
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
        )
      ) : (
        <EmptyState
          title="Connect Wallet"
          description="Connect your wallet to view transaction details."
          action={
            <div className="flex flex-col space-y-2">
              <ConnectWalletButton
                text={
                  ethereum
                    ? "Connect Wallet to continue"
                    : "Install MetaMask to continue"
                }
              />
              <Button variant={"link"} onClick={() => router.push("/personal")}>
                Go back
              </Button>
            </div>
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
