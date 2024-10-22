import React from "react";
import Layout from "../components/layouts/default-layout";
import Home from "../components/modules/home";
import { Metadata } from "next/types";
import { MarketPlace } from "@/pages/catalogs";

export const metadata: Metadata = {
  title: "Tenscan Gateway",
  description: "Tenscan Gateway: A gateway to the Tenscan ecosystem",
};

export default function HomePage() {
  return (
    <Layout>
      <div className="flex items-center justify-center w-full h-full">
        <Home />
      </div>

      <MarketPlace />
    </Layout>
  );
}
