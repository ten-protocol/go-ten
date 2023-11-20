import React from "react";
import { Metadata } from "next";
import Layout from "@/components/layouts/default-layout";
import Home from "@/components/modules/home";

export const metadata: Metadata = {
  title: "Tenscan Gateway",
  description: "Tenscan Gateway: A gateway to the Tenscan ecosystem",
};

export default function DashboardPage() {
  return (
    <Layout>
      <Home />
    </Layout>
  );
}
