import { fetchRollupByHash } from "@/api/rollups";
import Layout from "@/src/components/layouts/default-layout";
import { RollupDetailsComponent } from "@/src/components/modules/rollups/rollup-details";
import EmptyState from "@repo/ui/components/common/empty-state";
import LoadingState from "@repo/ui/components/common/loading-state";
import { Button } from "@repo/ui/components/shared/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@repo/ui/components/shared/card";
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

  return (
    <Layout>
      {isLoading ? (
        <LoadingState numberOfItems={10} />
      ) : rollupDetails ? (
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Rollup #{Number(rollupDetails?.ID)}</CardTitle>
            <CardDescription>
              Overview of rollup #{Number(rollupDetails?.ID)}
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

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
