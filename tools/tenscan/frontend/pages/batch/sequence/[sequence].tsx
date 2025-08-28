import { fetchBatchBySequence } from "@/api/batches";
import Layout from "@/src/components/layouts/default-layout";
import { BatchDetailsComponent } from "@/src/components/modules/batches/batch-details";
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

export default function BatchSequenceDetails() {
  const router = useRouter();
  const { sequence } = router.query;

  const { data, isLoading } = useQuery({
    queryKey: ["batchSequenceDetails", sequence],
    queryFn: () => fetchBatchBySequence(sequence as string),
    enabled: !!sequence,
  });

  const batchDetails = data?.item;

  return (
    <Layout>
      {isLoading ? (
        <LoadingState numberOfItems={10} />
      ) : batchDetails ? (
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Batch #{Number(sequence)}</CardTitle>
            <CardDescription>
              Overview of the Batch with sequence #{Number(sequence)}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <BatchDetailsComponent batchDetails={batchDetails} />
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
