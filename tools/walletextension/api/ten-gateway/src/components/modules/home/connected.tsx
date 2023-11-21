import { useWalletConnection } from "@/components/providers/wallet-provider";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import React from "react";

const Accounts = [
  {
    name: "Account 1",
    connected: true,
  },
  {
    name: "Account 2",
    connected: true,
  },
  {
    name: "Account 3",
    connected: true,
  },
  {
    name: "Account 4",
    connected: false,
  },
];

const Connected = () => {
  const { accounts } = useWalletConnection();

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Account</TableHead>
          <TableHead>Connected</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {Accounts.map((account, i) => (
          <TableRow key={i}>
            <TableCell className="font-medium">{account.name}</TableCell>
            <TableCell>
              <Badge variant={account.connected ? "success" : "destructive"}>
                {account.connected ? "Yes" : "No"}
              </Badge>
            </TableCell>
            <TableCell>
              <Button size={"sm"}>
                {account.connected ? "Disconnect" : "Connect"}
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};
export default Connected;
