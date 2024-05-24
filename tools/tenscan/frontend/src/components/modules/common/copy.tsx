import React from "react";
import { CopyIcon, CheckIcon } from "@radix-ui/react-icons";
import { useCopy } from "@/src/hooks/useCopy";
import { Button } from "@repo/ui/shared/button";

const Copy = ({ value }: { value: string | number }) => {
  const { copyToClipboard, copied } = useCopy();
  return (
    <Button
      type="submit"
      variant={"clear"}
      size="sm"
      className="text-muted-foreground"
      onClick={() => copyToClipboard(value.toString())}
    >
      <span className="sr-only">Copy</span>
      {copied ? <CheckIcon className="text-primary" /> : <CopyIcon />}
    </Button>
  );
};

export default Copy;
