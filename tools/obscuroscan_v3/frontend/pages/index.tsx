import React from "react";
import { Metadata } from "next";
import Layout from "@/src/components/layouts/default-layout";
import Dashboard from "@/src/components/modules/dashboard";

export const metadata: Metadata = {
  title: "Dashboard",
  description: "Obscuroscan Dashboard",
};

export default function DashboardPage() {
  return (
    <Layout>
      <Dashboard />
    </Layout>
  );
}
