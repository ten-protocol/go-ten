import React from "react";
import { fetchBatchTransactions } from "../../../api/batches";
import Layout from "../../../src/components/layouts/default-layout";
import { DataTable } from "@repo/ui/components/common/data-table/data-table";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { columns } from "../../../src/components/modules/batches/transaction-columns";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@repo/ui/components/shared/card";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";
import { getOptions } from "../../../src/lib/constants";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";

export default function BatchTransactions() {
  const router = useRouter();
  const { hash } = router.query;
  const options = getOptions(router.query);

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["batchTransactions", { hash, options }],
    queryFn: () => fetchBatchTransactions(hash as string, options),
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
          <div className="flex space-x-2 wrap">
            <CardDescription className="flex items-center space-x-2">
              <span>
                Overview of all transactions in this batch:
                <TruncatedAddress
                  address={hash as string}
                  showCopy={false}
                  link={pathToUrl(pageLinks.batchByHash, {
                    hash: hash as string,
                  })}
                />
              </span>
            </CardDescription>
          </div>
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
