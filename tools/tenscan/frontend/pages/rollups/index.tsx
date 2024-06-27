import React from "react";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { useRollupsService } from "@/src/services/useRollupsService";
import { Metadata } from "next";
import { formatNumber } from "@/src/lib/utils";
import { columns } from "@/src/components/modules/rollups/columns";

export const metadata: Metadata = {
  title: "Rollups",
  description: "A table of rollups.",
};

export default function Rollups() {
  const { rollups, setNoPolling, isRollupsLoading, refetchRollups } =
    useRollupsService();
  const { RollupsData, Total } = rollups?.result || {
    RollupsData: [],
    Total: 0,
  };

  React.useEffect(() => {
    setNoPolling(true);

    return () => {
      setNoPolling(false);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Rollups</h2>
            {/* uncomment the following line when total count feature is implemented */}
            {/* <p className="text-sm text-muted-foreground">
              {formatNumber(Total)} Rollups found.
            </p> */}
          </div>
        </div>
        {RollupsData ? (
          <DataTable
            columns={columns}
            data={RollupsData}
            refetch={refetchRollups}
            total={+Total}
            isLoading={isRollupsLoading}
          />
        ) : (
          <div>No rollups found.</div>
        )}
      </div>
    </Layout>
  );
}
