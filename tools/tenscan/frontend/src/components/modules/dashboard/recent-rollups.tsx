import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { formatTimeAgo } from "@repo/ui/lib/utils";
import { Rollup } from "@/src/types/interfaces/RollupInterfaces";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";
import GlitchTextAnimation from "@/src/components/GlitchTextAnimation";
import { RecentItemsList } from "./recent-items-list";

export function RecentRollups({ rollups }: { rollups: any }) {
  const renderRollupItem = (rollup: Rollup, isNewItem: boolean) => (
    <>
      <div className="ml-4 space-y-1 relative z-10">
        <p className="text-sm font-medium leading-none">
          <GlitchTextAnimation text={`#${Number(rollup?.ID)}`} hover={false} active={isNewItem} onView={false} />
        </p>
        <p className="text-sm text-muted-foreground word-break-all">
          <GlitchTextAnimation text={formatTimeAgo(rollup?.Timestamp)} hover={false} active={isNewItem} onView={false} />
        </p>
      </div>
      <div className="ml-auto font-medium min-w-[140px] relative z-10" onClick={(e) => e.stopPropagation()}>
        <TruncatedAddress
          address={rollup?.Hash}
          animate={isNewItem}
          AnimationComponent={GlitchTextAnimation}
          showPopover={false}
        />
      </div>
    </>
  );

  const headers = (
    <div className="flex items-center text-xs font-medium text-muted-foreground uppercase tracking-wide">
      <div className="ml-2 flex-1">
        <span>Rollup</span>
      </div>
      <div className="min-w-[140px] text-right mr-8">
        <span>Hash</span>
      </div>
    </div>
  );

  return (
    <RecentItemsList
      items={rollups?.result?.RollupsData || []}
      getItemId={(rollup: Rollup) => rollup.ID.toString()}
      getItemLink={(rollup: Rollup) => pathToUrl(pageLinks.rollupByHash, { hash: rollup?.Hash })}
      renderItem={renderRollupItem}
      headers={headers}
    />
  );
}