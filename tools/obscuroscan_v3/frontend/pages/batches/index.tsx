import React from "react";
import { columns } from "@/src/components/modules/batches/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import { useBatchesService } from "@/src/hooks/useBatchesService";

export const metadata: Metadata = {
  title: "Batches",
  description: "A table of Batches.",
};

export default function Batches() {
  const { batches } = useBatchesService();

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
          {batches?.result?.BatchesData ? (
            <DataTable columns={columns} data={batches?.result?.BatchesData} />
          ) : (
            <p>Loading...</p>
          )}
        </div>
      </Layout>
    </>
  );
}
