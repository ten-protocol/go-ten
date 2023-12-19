import React from "react";
import Header from "./header";
import Footer from "./footer";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <div className="bg-background">
      <div className="flex flex-col min-h-screen max-w-[1400px] mx-auto px-4">
        <Header />
        <div className="flex-1 space-y-4 py-6">{children}</div>
        <Footer />
      </div>
    </div>
  );
};

export default Layout;
