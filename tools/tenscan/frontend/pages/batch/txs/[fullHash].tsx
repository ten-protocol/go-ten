import { fetchBatchTransactions } from "@/api/batches";
import Layout from "@/src/components/layouts/default-layout";
import { DataTable } from "@repo/ui/common/data-table/data-table";
import EmptyState from "@repo/ui/common/empty-state";
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
      {TransactionsData ? (
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Transactions</CardTitle>
            <CardDescription>
              Overview of transactions at batch{" "}
              {"#" + TransactionsData[0]?.BatchHeight}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <DataTable
              columns={columns}
              data={TransactionsData}
              refetch={refetch}
              total={+Total}
              isLoading={isLoading}
            />
          </CardContent>
        </Card>
      ) : (
        <EmptyState
          title="No transactions found"
          description="There are no transactions in this batch."
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
