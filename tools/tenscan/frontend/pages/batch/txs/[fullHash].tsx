import { fetchBatchTransactions } from "@/api/batches";
import Layout from "@/src/components/layouts/default-layout";
import { DataTable } from "@repo/ui/common/data-table/data-table";
import { columns } from "@/src/components/modules/batches/transaction-columns";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@repo/ui/shared/card";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";
import { getOptions } from "@/src/lib/constants";
import TruncatedAddress from "@repo/ui/shared/truncated-address";

export default function BatchTransactions() {
  const router = useRouter();
  const { fullHash } = router.query;
  const options = getOptions(router.query);

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["batchTransactions", { fullHash, options }],
    queryFn: () => fetchBatchTransactions(fullHash as string, options),
  });

  const { TransactionsData, Total } = data?.result || {
    TransactionsData: [],
    Total: 0,
  };

  return (
    <Layout>
      <Card className="col-span-3">
        <CardHeader>
          <CardTitle>Transactions</CardTitle>
          <CardDescription className="flex items-center space-x-2">
            <p>Overview of all transactions in this batch:</p>
            <TruncatedAddress
              address={fullHash as string}
              showCopy={false}
              link={"/batch/" + fullHash}
            />
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DataTable
            columns={columns}
            data={TransactionsData}
            refetch={refetch}
            total={+Total}
            isLoading={isLoading}
            noResultsMessage="No transactions found in this batch."
            noPagination={true}
          />
        </CardContent>
      </Card>
    </Layout>
  );
}

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
