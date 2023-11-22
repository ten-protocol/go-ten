import React from "react";
import { Button } from "@/components/ui/button";
import { useCopy } from "@/hooks/useCopy";
import { CopyIcon, CheckIcon } from "@radix-ui/react-icons";

const Copy = ({ value }: { value: string | number }) => {
  const { copyToClipboard, copied } = useCopy();
  return (
    <Button
      type="submit"
      variant={"clear"}
      size="sm"
      className="px-3 py-1 text-muted-foreground"
      onClick={() => copyToClipboard(value.toString())}
    >
      <span className="sr-only">Copy</span>
      {copied ? <CheckIcon /> : <CopyIcon />}
    </Button>
  );
};

export default Copy;
