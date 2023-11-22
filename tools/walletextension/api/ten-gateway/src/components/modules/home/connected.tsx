import { useWalletConnection } from "@/components/providers/wallet-provider";
import { Badge } from "@/components/ui/badge";
import { Button, LinkButton } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Account } from "@/types/interfaces/WalletInterfaces";
import React from "react";
import TruncatedAddress from "../common/truncated-address";
import { socialLinks } from "@/lib/constants";
import { Skeleton } from "@/components/ui/skeleton";

const Connected = () => {
  const { accounts, connectAccount, disconnectAccount, revokeAccounts } =
    useWalletConnection();
  console.log("🚀 ~ file: connected.tsx:20 ~ Connected ~ accounts:", accounts);

  return (
    <>
      <div>
        <h1 className="text-4xl font-bold">Connected Accounts</h1>
        <h3 className="text-sm text-muted-foreground my-4">
          Manage the accounts you have connected to the Ten Gateway. You can
          revoke access to your accounts at any time and request new tokens from
          the Ten Discord.
        </h3>
        <div className="flex justify-end space-x-2 my-4">
          <LinkButton size={"sm"} href={socialLinks.discord}>
            Request Tokens
          </LinkButton>
          <Button size={"sm"} variant={"destructive"} onClick={revokeAccounts}>
            Revoke Accounts
          </Button>
        </div>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Account</TableHead>
            <TableHead>Connected</TableHead>
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {!accounts ? (
            <TableRow>
              <Skeleton className="w-full" />
            </TableRow>
          ) : accounts.length === 0 ? (
            <TableRow>
              <TableCell colSpan={3} className="text-center">
                No accounts connected
              </TableCell>
            </TableRow>
          ) : (
            accounts.map((account: Account, i: number) => (
              <TableRow key={i}>
                <TableCell className="font-medium">
                  <TruncatedAddress address={account.name} />
                </TableCell>
                <TableCell>
                  <Badge
                    variant={account.connected ? "success" : "destructive"}
                  >
                    {account.connected ? "Yes" : "No"}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Button
                    size={"sm"}
                    variant={account.connected ? "outline" : "default"}
                    onClick={
                      account.connected
                        ? () => disconnectAccount(account.name)
                        : () => connectAccount(account.name)
                    }
                  >
                    {account.connected ? "Disconnect" : "Connect"}
                  </Button>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </>
  );
};
export default Connected;
