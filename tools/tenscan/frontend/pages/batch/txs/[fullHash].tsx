import { fetchBatchTransactions } from "@/api/batches";
import Layout from "@/src/components/layouts/default-layout";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import EmptyState from "@/src/components/modules/common/empty-state";
import { columns } from "@/src/components/modules/batches/transaction-columns";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/src/components/ui/card";
import { formatNumber } from "@/src/lib/utils";
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
            <CardTitle>
              Batch Transactions at {"#" + TransactionsData[0]?.BatchHeight}
            </CardTitle>
            <CardDescription>
              Overview of the batch transactions at{" "}
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
