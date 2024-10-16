import React from "react";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import TransactionsComponent from "@/src/components/modules/transactions";

export const metadata: Metadata = {
  title: "Transactions",
  description: "A table of transactions.",
};

export default function Transactions() {
  return (
    <Layout>
      <TransactionsComponent />
    </Layout>
  );
}
