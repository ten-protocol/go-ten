import React from "react";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import PersonalTransactions from "@/src/components/modules/personal";

export const metadata: Metadata = {
  title: "Personal Transactions",
  description: "Tenscan Personal Transactions",
};

export default function PersonalPage() {
  return (
    <Layout>
      <PersonalTransactions />
    </Layout>
  );
}
