import React from "react";
import { columns } from "@/src/components/modules/blocks/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import { useBlocksService } from "@/src/services/useBlocksService";
import { firstItem, lastItem } from "@/src/lib/utils";

export const metadata: Metadata = {
  title: "Blocks",
  description: "A table of Blocks.",
};

export default function Blocks() {
  const { blocks, setNoPolling, refetchBlocks, isBlocksLoading } =
    useBlocksService();
  const { BlocksData, Total } = blocks?.result || {
    BlocksData: [],
    Total: 0,
  };

  React.useEffect(() => {
    setNoPolling(true);
    return () => setNoPolling(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const firstBlockNumber = Number(firstItem(BlocksData, "blockHeader.number"));
  const lastBlockNumber = Number(lastItem(BlocksData, "blockHeader.number"));

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Blocks</h2>
            {BlocksData?.length > 0 && (
              <p className="text-sm text-muted-foreground">
                Showing blocks #{firstBlockNumber}{" "}
                {lastBlockNumber !== firstBlockNumber &&
                  "to #" + lastBlockNumber}
                {/* uncomment the following line when total count feature is implemented */}
                {/* of {formatNumber(Total)} blocks. */}
              </p>
            )}
          </div>
        </div>
        {BlocksData ? (
          <DataTable
            columns={columns}
            data={BlocksData}
            total={+Total}
            refetch={refetchBlocks}
            isLoading={isBlocksLoading}
          />
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </Layout>
  );
}
