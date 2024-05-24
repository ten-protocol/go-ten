import { fetchBatchByHash } from "@/api/batches";
import Layout from "@/src/components/layouts/default-layout";
import { BatchHashDetailsComponent } from "@/src/components/modules/batches/batch-hash-details";
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

export default function Batch() {
  const router = useRouter();
  const { hash } = router.query;

  const { data, isLoading } = useQuery({
    queryKey: ["batch", hash],
    queryFn: () => fetchBatchByHash(hash as string),
  });

  const batchDetails = data?.item;

  return (
    <Layout>
      {isLoading ? (
        <Skeleton className="h-6 w-24" />
      ) : batchDetails ? (
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Batch #{Number(batchDetails?.Header?.number)}</CardTitle>
            <CardDescription>
              Overview of the batch #{Number(batchDetails?.Header?.number)}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <BatchHashDetailsComponent batchDetails={batchDetails} />
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
