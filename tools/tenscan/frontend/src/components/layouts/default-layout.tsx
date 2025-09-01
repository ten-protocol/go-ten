import React from "react";
import Header from "./header";
import Footer from "./footer";
import { Chakra_Petch } from 'next/font/google';

interface LayoutProps {
  children: React.ReactNode;
}

// const geistSans = Geist({
//     variable: '--font-geist-sans',
//     subsets: ['latin'],
// });
//
// const geistMono = Geist_Mono({
//     variable: '--font-geist-mono',
//     subsets: ['latin'],
// });

const chakra = Chakra_Petch({
    subsets: ['latin'],
    variable: '--font-chakra',
    display: 'swap',
    weight: ['300', '400', '500', '600', '700'],
});


const Layout = ({ children }: LayoutProps) => {
  return (
      <div className={`bg-background ${chakra.variable} antialiased`}>
          <div className="fixed inset-0 pointer-events-none z-40 opacity-[.03] grain-overlay"/>
          <div className="flex flex-col min-h-screen max-w-[1400px] mx-auto px-4">
              <Header/>
              <div className="flex-1 space-y-4 py-6">{children}</div>
              <Footer/>
          </div>
          <div className="bg-after-glow"/>
      </div>
  );
};

export default Layout;
