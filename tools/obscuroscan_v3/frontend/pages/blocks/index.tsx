import React from "react";
import { columns } from "@/src/components/modules/blocks/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import { useBlocksService } from "@/src/hooks/useBlocksService";

export const metadata: Metadata = {
  title: "Blocks",
  description: "A table of Blocks.",
};

export default function Blocks() {
  const { blocks } = useBlocksService();

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Blocks</h2>
            <p className="text-muted-foreground">A table of Blocks.</p>
          </div>
        </div>
        {blocks?.result?.BlocksData ? (
          <DataTable columns={columns} data={blocks?.result?.BlocksData} />
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </Layout>
  );
}
