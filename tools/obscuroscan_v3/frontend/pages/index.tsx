import React from "react";
import { Metadata } from "next";
import Layout from "@/components/layouts/default-layout";
import Dashboard from "@/components/modules/dashboard";

export const metadata: Metadata = {
  title: "Dashboard",
  description: "Example dashboard app built using the components.",
};

export default function DashboardPage() {
  return (
    <Layout>
      <Dashboard />
    </Layout>
  );
}
