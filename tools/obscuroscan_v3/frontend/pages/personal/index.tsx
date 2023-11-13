import React from "react";
import { Metadata } from "next";
import Layout from "@/components/layouts/default-layout";
import Personal from "@/components/modules/dashboard";

export const metadata: Metadata = {
  title: "Personal Transactions",
  description: "ObscuroScan Personal Transactions",
};

export default function DashboardPage() {
  return (
    <Layout>
      <Personal />
    </Layout>
  );
}
