// {
//     "Header":
//     {
//     "parentHash":
//     "0x558b1be28fe24a6766e40aa24317492fe5978b159f8fd543170d61ae2749a8bb",
//     "stateRoot":
//     "0xb6c8ea4a0cfc96e202b512a2096c313e36c15883227b83f146564cdea5484e82",
//     "transactionsRoot":
//     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
//     "receiptsRoot":
//     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
//     "number":
//     9084,
//     "sequencerOrderNo":
//     9085,
//     "gasLimit":
//     1537228672809129200,
//     "gasUsed":
//     0,
//     "timestamp":
//     1699909432,
//     "extraData":
//     "",
//     "baseFee":
//     1,
//     "coinbase":
//     "0xd6c9230053f45f873cb66d8a02439380a37a4fbf",
//     "l1Proof":
//     "0x712ca13aff6745094fc215ca9a6b9790e778f13f852d0439b6f998e0b49b64c2",
//     "R":
//     1.0666142563506898e+77,
//     "S":
//     1.9832304678215498e+76,
//     "crossChainMessages":
//     [
//     ],
//     "inboundCrossChainHash":
//     "0x712ca13aff6745094fc215ca9a6b9790e778f13f852d0439b6f998e0b49b64c2",
//     "inboundCrossChainHeight":
//     37008,
//     "transfersTree":
//     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
//     "hash":
//     "0xa6cecb66607c02561dce716d6c077cfaea40119b2ad30427474801cea3740d4a",
//     "sha3Uncles":
//     null,
//     "miner":
//     "0xd6c9230053f45f873cb66d8a02439380a37a4fbf",
//     "logsBloom":
//     null,
//     "difficulty":
//     null,
//     "nonce":
//     null,
//     "baseFeePerGas":
//     1
//     },
//     "TxHashes":
//     [
//     ],
//     "EncryptedTxBlob":
//     "Fse8O1ZX32W9p68bd8ExeNiPMvfHnNi90o8pgFCIjiQB"

import { Separator } from "@/components/ui/separator";
import TruncatedAddress from "../common/truncated-address";

//     }
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
