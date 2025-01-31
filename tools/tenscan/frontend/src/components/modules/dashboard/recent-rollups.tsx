import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { formatTimeAgo } from "@repo/ui/lib/utils";
import { Avatar, AvatarFallback } from "@repo/ui/components/shared/avatar";
import { Rollup } from "@/src/types/interfaces/RollupInterfaces";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";

export function RecentRollups({ rollups }: { rollups: any }) {
  return (
    <div className="space-y-8">
      {rollups?.result?.RollupsData?.map((rollup: Rollup, i: number) => (
        <div className="flex items-center" key={i}>
          <Avatar className="h-9 w-9">
            <AvatarFallback>RP</AvatarFallback>
          </Avatar>
          <div className="ml-4 space-y-1">
            <p className="text-sm font-medium leading-none">
              #{Number(rollup?.ID)}
            </p>
            <p className="text-sm text-muted-foreground word-break-all">
              {formatTimeAgo(rollup?.Timestamp)}
            </p>
          </div>
          <div className="ml-auto font-medium min-w-[140px]">
            <TruncatedAddress
              address={rollup?.Hash}
              link={pathToUrl(pageLinks.rollupByHash, { hash: rollup?.Hash })}
            />
          </div>
        </div>
      ))}
    </div>
  );
}
