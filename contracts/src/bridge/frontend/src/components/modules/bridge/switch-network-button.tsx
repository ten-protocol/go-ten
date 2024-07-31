import { ArrowDownUpIcon, Loader } from "lucide-react";
import { Button } from "../../ui/button";

export const SwitchNetworkButton = ({ handleSwitchNetwork, loading }: any) => {
  return (
    <div className="flex items-center justify-center">
      <Button
        type="button"
        className="mt-4"
        variant="outline"
        size="sm"
        onClick={handleSwitchNetwork}
        disabled={loading}
      >
        <ArrowDownUpIcon className="h-4 w-4" />
      </Button>
    </div>
  );
};
