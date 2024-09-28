import React, { useEffect, useState } from "react";
import {
  Alert,
  AlertTitle,
  AlertDescription,
} from "@repo/ui/components/shared/alert";
import { Badge } from "@repo/ui/components/shared/badge";
import { Button } from "@repo/ui/components/shared/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@repo/ui/components/shared/card";
import { Textarea } from "@repo/ui/components/shared/textarea";
import { currentEncryptedKey } from "@/src/lib/constants";
import { CopyIcon, Terminal } from "@repo/ui/components/shared/react-icons";
import { useRouter } from "next/router";
import JSONPretty from "react-json-pretty";
import { useRollupsService } from "@/src/services/useRollupsService";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@repo/ui/components/shared/tooltip";

export default function Decrypt() {
  const router = useRouter();
  const { decryptedRollup, decryptEncryptedData } = useRollupsService();

  const [encryptedRollup, setEncryptedRollup] = useState<string>("");

  useEffect(() => {
    if (router.query.encryptedString) {
      setEncryptedRollup(
        decodeURIComponent(router.query.encryptedString as string)
      );
    }
  }, [router.query.encryptedString]);

  const decrypt = (event: any) => {
    event.preventDefault();
    decryptEncryptedData({ StrData: encryptedRollup });
    setEncryptedRollup("");
  };

  const handleTextareaChange = (event: any) => {
    setEncryptedRollup(event.target.value);
  };

  return (
    <>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle>Static Keys</CardTitle>
        </CardHeader>
        <CardContent>
          <Alert className="mt-2" variant={"warning"}>
            <Terminal className="h-4 w-4" />
            <AlertTitle>
              Decrypting transaction blobs is only possible on testnet!
            </AlertTitle>
            <AlertDescription>
              The rollup encryption key is long-lived and well-known. On
              mainnet, rollups will use rotating keys that are not known to
              anyone - or anything - other than the Obscuro enclaves.
            </AlertDescription>
          </Alert>
          <div className="mt-2">
            <Badge variant={"default"}>Current Static Encryption Key</Badge>
            <div className="flex items-center space-x-2 mt-2">
              <div className="truncate max-w-[700px]">
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger>
                      <p className="text-sm font-medium truncate">
                        {currentEncryptedKey}
                      </p>
                    </TooltipTrigger>
                    <TooltipContent>{currentEncryptedKey}</TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </div>
              <CopyIcon className="h-4 w-4 ml-2" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="mt-4">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle>Encrypted Rollup</CardTitle>
        </CardHeader>
        <CardContent className="mt-2">
          <form onSubmit={decrypt}>
            <Textarea
              placeholder="Please enter Encrypted Rollup"
              onChange={handleTextareaChange}
              value={encryptedRollup}
            />
            <Button className="mt-4">Decrypt</Button>
          </form>
        </CardContent>
      </Card>
      {decryptedRollup ? (
        <JSONPretty id="json-pretty" data={decryptedRollup}></JSONPretty>
      ) : null}
    </>
  );
}
