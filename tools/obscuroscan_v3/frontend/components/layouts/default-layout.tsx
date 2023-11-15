import React from "react";
import Header from "./header";
import Footer from "./footer";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <div className="flex-col md:flex">
      <div className="flex flex-col min-h-screen min-w-[1400px] mx-auto">
        <Header />
        <div className="flex-1 space-y-4 p-8 pt-6">{children}</div>
        <Footer />
      </div>
    </div>
  );
};

export default Layout;
