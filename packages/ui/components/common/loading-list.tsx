import KeyValueItem, { KeyValueList } from "../shared/key-value";
import { Skeleton } from "../shared/skeleton";

export function LoadingList({ numberOfItems = 4 }: { numberOfItems?: number }) {
  return (
    <div className="space-y-8">
      <KeyValueList>
        {Array.from({ length: numberOfItems }).map((_, index) => (
          <KeyValueItem
            key={index}
            label={<Skeleton className="h-6 w-24" />}
            value={<Skeleton className="h-6 w-24" />}
            isLastItem={index === numberOfItems - 1}
          />
        ))}
      </KeyValueList>
    </div>
  );
}
