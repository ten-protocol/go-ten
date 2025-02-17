import { fetchBatchByHeight } from "@/api/batches";
import Layout from "@/src/components/layouts/default-layout";
import { BatchHeightDetailsComponent } from "@/src/components/modules/batches/batch-height-details";
import LoadingState from "@repo/ui/components/common/loading-state";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@repo/ui/components/shared/card";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";

export default function Batch() {
  const router = useRouter();
  const { height } = router.query;

  const { data, isLoading } = useQuery({
    queryKey: ["batchHeight", height],
    queryFn: () => fetchBatchByHeight(height as string),
  });

  const batchDetails = data?.item;

  return (
    <Layout>
      {isLoading ? (
        <LoadingState numberOfItems={10} />
      ) : batchDetails ? (
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Batch #{Number(batchDetails?.header?.number)}</CardTitle>
            <CardDescription>
              Overview of the batch #{Number(batchDetails?.header?.number)}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <BatchHeightDetailsComponent batchDetails={batchDetails} />
          </CardContent>
        </Card>
      ) : (
        <div>Batch not found</div>
      )}
    </Layout>
  );
}

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
