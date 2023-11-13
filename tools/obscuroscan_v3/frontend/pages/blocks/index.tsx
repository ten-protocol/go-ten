import React from "react";
import { columns } from "@/components/modules/blocks/columns";
import { DataTable } from "@/components/modules/common/data-table/data-table";
import Layout from "@/components/layouts/default-layout";
import { Metadata } from "next";
import { useBlocks } from "@/src/hooks/useBlocks";

export const metadata: Metadata = {
  title: "Blocks",
  description: "A table of Blocks.",
};

export default function Blocks() {
  const { blocks } = useBlocks();

  return (
    <>
      <Layout>
        <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
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
    </>
  );
}
