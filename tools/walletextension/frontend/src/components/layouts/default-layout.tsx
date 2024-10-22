import React from "react";
import Header from "./header";
import Footer from "./footer";
import { SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/layouts/app-sidebar";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <div
      className={
        "bg-background md:min-w-[450px] min-h-screen max-w-[1400px] mx-auto px-4"
      }
    >
      <SidebarProvider>
        <AppSidebar />
        <main className={"w-full"}>
          <Header />
          <div className="flex-1 space-y-4 py-6">{children}</div>
          <Footer />
        </main>
      </SidebarProvider>
    </div>
  );
};

export default Layout;
