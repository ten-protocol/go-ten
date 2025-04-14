'use client';
import { useAccount } from 'wagmi';
import WalletConnected from "@/components/WalletConnected/WalletConnected";
import DisconnectedWallet from "@/components/DisconnectedWallet/DisconnectedWallet";
import Header from "@/components/Header/Header";
import Footer from "@/components/Footer/Footer";
import PromoApps from "@/components/PromoApps/PromoApps";

export default function Home() {
  const {isConnected } = useAccount();

  return (
      <div className="min-h-screen flex flex-col items-center justify-center p-8 overflow-y-hidden relative">
          <div className="fixed inset-0 pointer-events-none z-40 opacity-[.03] grain-overlay"/>
          <Header/>

          <main className="flex flex-col items-center justify-center flex-1 gap-8 mt-16">
              <div className="text-center mb-12">
                  <h1 className="text-[3rem] font-bold -mb-1">Welcome to the TEN Gateway!</h1>
                  <h2 className="opacity-80 text-lg">
                      Your portal into the universe of encrypted Ethereum on TEN Protocol.
                  </h2>
              </div>

              {isConnected ? (
                  <WalletConnected/>
              ) : (
                  <DisconnectedWallet/>
              )}
              <PromoApps/>
          </main>


          <Footer/>
          <div className="bg-after-glow"/>

      </div>
  );
}
