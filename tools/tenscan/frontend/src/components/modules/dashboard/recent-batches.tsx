import TruncatedAddress from "../common/truncated-address";
import { formatTimeAgo } from "@/src/lib/utils";
import { Batch } from "@/src/types/interfaces/BatchInterfaces";
import { Avatar, AvatarFallback } from "@/src/components/ui/avatar";

export function RecentBatches({ batches }: { batches: any }) {
  return (
    <div className="space-y-8">
      {batches?.result?.BatchesData.map((batch: Batch, i: number) => (
        <div className="flex items-center" key={i}>
          <Avatar className="h-9 w-9">
            <AvatarFallback>BN</AvatarFallback>
          </Avatar>
          <div className="ml-4 space-y-1">
            <p className="text-sm font-medium leading-none">
              #{Number(batch?.number)}
            </p>
            <p className="text-sm text-muted-foreground word-break-all">
              {formatTimeAgo(batch?.timestamp)}
            </p>
          </div>
          <div className="ml-auto font-medium min-w-[140px]">
            <TruncatedAddress address={batch?.hash} />
          </div>
        </div>
      ))}
    </div>
  );
}
