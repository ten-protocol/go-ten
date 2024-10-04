import React from "react";
import { Button } from "../../ui/button";
import { useCopy } from "../../../hooks/useCopy";
import { CopyIcon, CheckIcon } from "@radix-ui/react-icons";

const Copy = ({ value }: { value: string | number | undefined }) => {
  const { copyToClipboard, copied } = useCopy();
  if (!value) return null;
  return (
    <Button
      type="button"
      variant={"clear"}
      size="sm"
      className="px-3 py-1 text-muted-foreground"
      onClick={() => (value ? copyToClipboard(String(value)) : null)}
    >
      <span className="sr-only">Copy</span>
      {copied ? <CheckIcon /> : <CopyIcon />}
    </Button>
  );
};

export default Copy;
