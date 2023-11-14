import React from "react";
import Layout from "@/components/layouts/default-layout";
import { Metadata } from "next";
import PersonalTransactions from "@/components/modules/personal";

export const metadata: Metadata = {
  title: "Personal Transactions",
  description: "ObscuroScan Personal Transactions",
};

export default function PersonalPage() {
  return (
    <Layout>
      <PersonalTransactions />
    </Layout>
  );
}
