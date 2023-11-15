import {
  TableCaption,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/src/components/ui/table";
import { Table } from "@/src/components/ui/table";
import { useContracts } from "@/src/hooks/useContracts";
import TruncatedAddress from "../../common/truncated-address";
import { Badge } from "@/src/components/ui/badge";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@/src/components/ui/card";

export default function VerifiedSequencerData() {
  const { sequencerData } = useContracts();

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle>Verified Sequencer Data</CardTitle>
      </CardHeader>
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
