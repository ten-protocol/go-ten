import { fetchRollupByHash } from "@/api/rollups";
import Layout from "@/src/components/layouts/default-layout";
import { RollupDetailsComponent } from "@/src/components/modules/rollups/rollup-details";
import EmptyState from "@/src/components/modules/common/empty-state";
import { Button } from "@/src/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/src/components/ui/card";
import { Skeleton } from "@/src/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";

export default function RollupDetails() {
  const router = useRouter();
  const { hash } = router.query;

  const { data, isLoading } = useQuery({
    queryKey: ["rollupDetails", hash],
    queryFn: () => fetchRollupByHash(hash as string),
  });

  const rollupDetails = data?.item;
  console.log("🚀 ~ RollupDetails ~ rollupDetails:", rollupDetails);

  return (
    <Layout>
      {isLoading ? (
        <Skeleton className="h-full w-full" />
      ) : rollupDetails ? (
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Rollup #{Number(rollupDetails?.ID)}</CardTitle>
            <CardDescription>
              Overview of the Rollup #{Number(rollupDetails?.ID)}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <RollupDetailsComponent rollupDetails={rollupDetails} />
          </CardContent>
        </Card>
      ) : (
        <EmptyState
          title="Rollup not found"
          description="The rollup you are looking for does not exist."
          action={
            <Button onClick={() => router.push("/rollups")}>Go back</Button>
          }
        />
      )}
    </Layout>
  );
}
