import React from "react";
import { Button } from "../../ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "../../ui/card";
import { Terminal } from "lucide-react";
import { useWalletConnection } from "../../providers/wallet-provider";
import Connected from "./connected";
import Disconnected from "./disconnected";

const Home = () => {
  const { walletConnected } = useWalletConnection();

  return (
    <div className="w-[800px] mx-auto">
      {walletConnected ? <Connected /> : <Disconnected />}
    </div>
  );
};

export default Home;
