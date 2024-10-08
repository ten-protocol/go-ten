import React from "react";
import Layout from "../../../src/components/layouts/default-layout";
import EmptyState from "@repo/ui/components/common/empty-state";
import { Button } from "@repo/ui/components/shared/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@repo/ui/components/shared/card";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";
import { fetchPersonalTxnByHash } from "../../../api/transactions";
import { PersonalTxnDetailsComponent } from "../../../src/components/modules/personal/personal-txn-details";
import { ethereum } from "@repo/ui/lib/utils";
import useWalletStore from "@repo/ui/stores/wallet-store";
import ConnectWalletButton from "@repo/ui/components/common/connect-wallet";
import LoadingState from "@repo/ui/components/common/loading-state";

export default function TransactionDetails() {
  const router = useRouter();
  const { provider, walletConnected } = useWalletStore();
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
          <LoadingState numberOfItems={10} />
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
