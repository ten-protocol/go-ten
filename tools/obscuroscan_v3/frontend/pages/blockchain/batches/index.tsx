import React from "react";
import { columns } from "@/components/modules/personal/batches/components/columns";
import { DataTable } from "@/components/modules/common/data-table/data-table";
import Layout from "@/components/layouts/default-layout";
import { Metadata } from "next";
import { useBatches } from "@/src/hooks/useBatches";

export const metadata: Metadata = {
  title: "Batches",
  description: "A table of Batches.",
};

export default function Batches() {
  const { batches } = useBatches();

  return (
    <>
      <Layout>
        <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
          <div className="flex items-center justify-between space-y-2">
            <div>
              <h2 className="text-2xl font-bold tracking-tight">Batches</h2>
              <p className="text-muted-foreground">A table of Batches.</p>
            </div>
          </div>
          {batches?.result?.batchesData ? (
            <DataTable columns={columns} data={batches?.result?.batchesData} />
          ) : (
            <p>Loading...</p>
          )}
        </div>
      </Layout>
    </>
  );
}
