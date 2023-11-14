import { Button } from "@/components/ui/button";
import TruncatedAddress from "../common/truncated-address";
import { formatTimeAgo } from "@/src/lib/utils";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { EyeOpenIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { Block } from "@/src/types/interfaces/BlockInterfaces";

export function RecentBlocks({ blocks }: { blocks: any }) {
  return (
    <div className="space-y-8">
      {blocks?.result?.BlocksData.map((block: Block, i: number) => (
        <div className="flex items-center" key={i}>
          <Avatar className="h-9 w-9">
            <AvatarFallback>BK</AvatarFallback>
          </Avatar>
          <div className="ml-4 space-y-1">
            <p className="text-sm font-medium leading-none">
              #{block?.blockHeader?.number}
            </p>
            <p className="text-sm text-muted-foreground word-break-all">
              {formatTimeAgo(block?.blockHeader?.timestamp)}
            </p>
          </div>
          <div className="ml-auto font-medium min-w-[140px]">
            <TruncatedAddress address={block?.blockHeader?.hash} />
          </div>
          <div className="ml-auto font-medium max-w-[50px]">
            <Link
              href={{
                pathname: `/blocks/${block?.blockHeader?.number}`,
              }}
            >
              <Button variant="link" size="sm">
                <EyeOpenIcon />
              </Button>
            </Link>
          </div>
        </div>
      ))}
    </div>
  );
}
