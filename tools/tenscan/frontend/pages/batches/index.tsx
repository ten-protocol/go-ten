import React from "react";
import { columns } from "@/src/components/modules/batches/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import { useBatchesService } from "@/src/services/useBatchesService";
import { firstItem, lastItem } from "@/src/lib/utils";

export const metadata: Metadata = {
  title: "Batches",
  description: "A table of Batches.",
};

export default function Batches() {
  const { batches, refetchBatches, isBatchesLoading, setNoPolling } =
    useBatchesService();
  const { BatchesData, Total } = batches?.result || {
    BatchesData: [],
    Total: 0,
  };

  React.useEffect(() => {
    setNoPolling(true);

    return () => {
      setNoPolling(false);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const firstBatchHeight = Number(firstItem(BatchesData, "height"));
  const lastBatchHeight = Number(lastItem(BatchesData, "height"));

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Batches</h2>
            {BatchesData?.length > 0 && (
              <p className="text-sm text-muted-foreground">
                Showing batches #{firstBatchHeight}{" "}
                {lastBatchHeight !== firstBatchHeight &&
                  "to #" + lastBatchHeight}
                {/* uncomment the following line when total count feature is implemented */}
                {/* of {formatNumber(Total)} batches. */}
              </p>
            )}
          </div>
        </div>
        {BatchesData ? (
          <DataTable
            columns={columns}
            data={BatchesData}
            refetch={refetchBatches}
            total={+Total}
            isLoading={isBatchesLoading}
          />
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </Layout>
  );
}
