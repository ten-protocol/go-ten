import {
  TableCaption,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@repo/ui/components/shared/table";
import { Table } from "@repo/ui/components/shared/table";
import { useContractsService } from "@/src/services/useContractsService";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { Badge } from "@repo/ui/components/shared/badge";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@repo/ui/components/shared/card";
import { Separator } from "@repo/ui/components/shared/separator";

export default function VerifiedSequencerData() {
  const { sequencerData } = useContractsService();

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle>Verified Sequencer Data</CardTitle>
      </CardHeader>
      <Separator />
      <CardContent>
        <Table>
          <TableCaption>Verified Sequencer Data</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>Contract Name</TableHead>
              <TableHead>Confirmed</TableHead>
              <TableHead>Contract Address</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sequencerData.map((contract, i) => (
              <TableRow key={i}>
                <TableCell className="font-medium">{contract.name}</TableCell>
                <TableCell>
                  <Badge
                    variant={contract.confirmed ? "success" : "destructive"}
                  >
                    {contract.confirmed ? "Yes" : "No"}
                  </Badge>
                </TableCell>
                <TableCell>
                  <TruncatedAddress address={contract.address} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
