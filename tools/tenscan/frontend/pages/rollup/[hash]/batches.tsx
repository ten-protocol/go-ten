import React from "react";
import { columns } from "@/src/components/modules/batches/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import { formatNumber } from "@/src/lib/utils";
import { useRollupsService } from "@/src/services/useRollupsService";

export const metadata: Metadata = {
  title: "Batches",
  description: "A table of Batches.",
};

export default function RollupBatches() {
  const { rollupBatches, isRollupBatchesLoading, refetchRollupBatches } =
    useRollupsService();

  const { BatchesData, Total } = rollupBatches?.result || {
    BatchesData: [],
    Total: 0,
  };

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Batches</h2>
            {/* comment out the following line when feature is implemented */}
            {/* <p className="text-sm text-muted-foreground">
              {formatNumber(Total)} Batch(es) found in this rollup.
            </p> */}
          </div>
        </div>
        {BatchesData ? (
          <DataTable
            columns={columns}
            data={BatchesData}
            refetch={refetchRollupBatches}
            total={+Total}
            isLoading={isRollupBatchesLoading}
          />
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </Layout>
  );
}

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
