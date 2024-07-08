import React from "react";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { useRollupsService } from "@/src/services/useRollupsService";
import { Metadata } from "next";
import { columns } from "@/src/components/modules/rollups/columns";
import { getItem } from "@/src/lib/utils";
import { ItemPosition } from "@/src/types/interfaces";

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

  const firstRollupID = Number(getItem(RollupsData, "ID"));
  const lastRollupID = Number(getItem(RollupsData, "ID", ItemPosition.LAST));

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Rollups</h2>
            {RollupsData?.length > 0 && (
              <p className="text-sm text-muted-foreground">
                Showing rollups #{firstRollupID}{" "}
                {lastRollupID !== firstRollupID && "to #" + lastRollupID}
                {/* uncomment the following line when total count feature is implemented */}
                {/* of {formatNumber(Total)} rollups. */}
              </p>
            )}
          </div>
        </div>
        <DataTable
          columns={columns}
          data={RollupsData}
          refetch={refetchRollups}
          total={+Total}
          isLoading={isRollupsLoading}
          noResultsWord="rollups"
        />
      </div>
    </Layout>
  );
}
