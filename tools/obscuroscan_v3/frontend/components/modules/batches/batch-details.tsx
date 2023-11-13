import { Separator } from "@/components/ui/separator";
import TruncatedAddress from "../common/truncated-address";

export function BatchDetails() {
  return (
    <div className="space-y-8">
      <div className="flex items-center">
        <div className="ml-4 space-y-1">
          <p className="text-sm font-medium leading-none">
            Batch Height: #9084
          </p>
          <p className="text-sm text-muted-foreground">
            <TruncatedAddress
              address={
                "0xa6cecb66607c02561dce716d6c077cfaea40119b2ad30427474801cea3740d4a"
              }
            />
          </p>
        </div>
      </div>
      <Separator />
      <div className="flex items-center">
        <div className="ml-4">
          <p>
            Parent Hash:
            0x558b1be28fe24a6766e40aa24317492fe5978b159f8fd543170d61ae2749a8bb
          </p>
        </div>
      </div>
    </div>
  );
}
