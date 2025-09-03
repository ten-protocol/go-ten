import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { formatTimeAgo } from "@repo/ui/lib/utils";
import { Batch } from "@/src/types/interfaces/BatchInterfaces";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";
import GlitchTextAnimation from "@/src/components/GlitchTextAnimation";
import { RecentItemsList } from "./recent-items-list";

export function RecentBatches({ batches }: { batches: any }) {
  const renderBatchItem = (batch: Batch, isNewItem: boolean) => (
    <>
      <div className="ml-4 space-y-1 relative z-10">
        <p className="text-sm font-medium leading-none">
          <GlitchTextAnimation text={`#${Number(batch?.height)}`} hover={false} active={isNewItem} onView={false} />
        </p>
        <p className="text-sm text-muted-foreground word-break-all">
          <GlitchTextAnimation text={formatTimeAgo(batch?.header?.timestamp)} hover={false} active={isNewItem} onView={false} />
        </p>
      </div>
      
      <div className="ml-auto font-medium min-w-[140px] relative z-10" onClick={(e) => e.stopPropagation()}>
        <TruncatedAddress
          address={batch?.header?.hash}
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
        <span>Batch</span>
      </div>
      <div className="min-w-[140px] text-right mr-8">
        <span>Hash</span>
      </div>
    </div>
  );

  return (
    <RecentItemsList
      items={batches?.result?.BatchesData || []}
      getItemId={(batch: Batch) => batch.header?.hash || batch.height?.toString() || ''}
      getItemLink={(batch: Batch) => pathToUrl(pageLinks.batchByHeight, { height: batch?.height.toString() })}
      renderItem={renderBatchItem}
      headers={headers}
    />
  );
}