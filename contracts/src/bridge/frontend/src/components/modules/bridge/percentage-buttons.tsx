import { PERCENTAGES } from "@/src/lib/constants";
import { Button } from "../../ui/button";

export const PercentageButtons = ({ setAmount }: any) => {
  return (
    <div className="flex items-center p-2">
      <div className="flex items-center space-x-2">
        {PERCENTAGES.map((percentage) => (
          <Button
            type="button"
            key={percentage.name}
            variant="outline"
            size="sm"
            className="dark:bg-[#292929]"
            onClick={() => setAmount(percentage.value)}
          >
            {percentage.name}
          </Button>
        ))}
      </div>
    </div>
  );
};
