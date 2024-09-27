import { Skeleton } from "../shared/skeleton";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "../shared/card";
import { LoadingList } from "./loading-list";

export default function LoadingState({
  numberOfItems = 4,
}: {
  numberOfItems?: number;
}) {
  return (
    <>
      <Card className="col-span-3">
        <CardHeader>
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-6 w-40" />
        </CardHeader>
        <CardContent>
          <LoadingList numberOfItems={numberOfItems} />
        </CardContent>
      </Card>
    </>
  );
}

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
