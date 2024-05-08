import React from "react";
import { Metadata } from "next";
import Layout from "@/src/components/layouts/default-layout";
import Dashboard from "@/src/components/modules/bridge";

export const metadata: Metadata = {
  title: "Ten Bridge",
  description: "Ten Bridge Dashboard",
};

export default function DashboardPage() {
  return (
    <Layout>
      <Dashboard />
    </Layout>
  );
}
