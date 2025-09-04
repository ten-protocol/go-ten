import { CopyIcon, CheckIcon } from "../shared/react-icons";
import { useCopy } from "../../hooks/useCopy";
import { Button } from "../shared/button";

const Copy = ({ value }: { value: string | number }) => {
  const { copyToClipboard, copied } = useCopy();
  return (
    <Button
      type="submit"
      variant={"clear"}
      size="sm"
      className="text-muted-foreground hover:text-primary"
      onClick={(e) => {
        e.stopPropagation();
        e.preventDefault();
        copyToClipboard(value.toString());
      }}
    >
      <span className="sr-only">Copy</span>
      {copied ? <CheckIcon className="text-primary" /> : <CopyIcon />}
    </Button>
  );
};

export default Copy;
