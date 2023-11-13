import { getBatchByHash } from "@/api/batches";
import Layout from "@/components/layouts/default-layout";
import { BatchDetails } from "@/components/modules/batches/batch-details";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";

export default function Batch() {
  const router = useRouter();
  const { batch } = router.query;

  const { data, isLoading } = useQuery({
    queryKey: ["Batch"],
    queryFn: () => getBatchByHash(batch as string),
  });

  return (
    <Layout>
      <Card className="col-span-3">
        <CardHeader>
          <CardTitle>Batch #{batch}</CardTitle>
          <CardDescription>
            {isLoading ? (
              <Skeleton className="h-6 w-24" />
            ) : (
              data?.result?.l1Proof
            )}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <BatchDetails />
        </CardContent>
      </Card>
    </Layout>
  );
}
