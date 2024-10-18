import React from "react";
import useGatewayService from "../../../services/useGatewayService";
import useWalletStore from "@/stores/wallet-store";
import { downloadMetaMask, ethereum } from "@repo/ui/lib/utils";
import { ConnectWalletButtonProps } from "@repo/ui/lib/interfaces/ui";
import {
  Link2Icon,
  LinkBreak2Icon,
} from "@repo/ui/components/shared/react-icons";
import ConnectWalletButton from "@repo/ui/components/common/connect-wallet";

const GatewayConnectButton = (
  props: Omit<ConnectWalletButtonProps, "onConnect" | "renderContent">
) => {
  const { walletConnected, revokeAccounts, accounts } = useWalletStore();
  const { connectToTenTestnet } = useGatewayService();

  const handleConnect = () => {
    if (ethereum) {
      if (walletConnected) {
        revokeAccounts();
      } else {
        connectToTenTestnet();
      }
    } else {
      downloadMetaMask();
    }
  };

  const renderGatewayContent = () => {
    if (!ethereum) {
      return (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          Install MetaMask
        </>
      );
    }

    if (walletConnected) {
      const connectedAccounts = accounts?.filter((acc) => acc.connected) || [];
      const totalAccounts = accounts?.length || 0;
      return (
        <>
          <LinkBreak2Icon className="h-4 w-4 mr-1" />
          {`${connectedAccounts.length}/${totalAccounts} Connected`}
        </>
      );
    } else {
      return (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          Connect to TEN Testnet
        </>
      );
    }
  };

  return (
    <ConnectWalletButton
      {...props}
      onConnect={handleConnect}
      renderContent={renderGatewayContent}
    />
  );
};

export default GatewayConnectButton;
