import { useWalletStore } from "../../providers/wallet-provider";
import { Button } from "../../ui/button";
import useGatewayService from "../../../services/useGatewayService";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";
import detectEthereumProvider from "@metamask/detect-provider";
import { ethereum, trackEvent } from "@/src/lib/utils";
import { toast } from "../../ui/use-toast";
import Web3Service from "@/src/services/web3service";
import { l1Bridge } from "@/src/lib/constants";
const ConnectWalletButton = () => {
  const { provider, signer, address, setProvider, setAddress } =
    useWalletStore();
  async function connectMetamask() {
    const provider = await detectEthereumProvider();

    if (provider) {
      // @ts-ignore
      const chainId = await provider.request({ method: "eth_chainId" });
      if (chainId !== "0x1bb") {
        return toast({
          title: "Wrong Network",
          description:
            'Not connected to Ten ! Connect at <a href="https://testnet.ten.xyz/" target="_blank" rel="noopener noreferrer">https://testnet.ten.xyz/</a> ',
          variant: "info",
        });
      }

      // Request account access if needed
      // @ts-ignore
      const accounts = await provider.request({
        method: "eth_requestAccounts",
      });

      // Set provider and address in the store
      setProvider(provider);
      setAddress(accounts[0]);

      trackEvent("connect_wallet", {
        value: accounts[0],
      });

      toast({
        title: "Connected",
        description: "Connected to wallet ! Account: " + accounts[0],
        variant: "success",
      });

      new Web3Service(signer);
    } else {
      console.error("Please install MetaMask!");
      toast({
        title: "MetaMask not found",
        description: "Please install MetaMask!",
        variant: "destructive",
      });
    }
  }
  return (
    <Button
      className="text-sm font-medium leading-none"
      variant={"outline"}
      // onClick={
      //   ethereum
      //     ? walletConnected
      //       ? revokeAccounts
      //       : connectToTenTestnet
      //     : downloadMetaMask
      // }
      // suppressHydrationWarning
    >
      {/* {walletConnected ? (
        <>
          <LinkBreak2Icon className="h-4 w-4 mr-1" />
          Disconnect
        </>
      ) : ( */}
      <>
        <Link2Icon className="h-4 w-4 mr-1" />
        {/* {ethereum ? "Connect" : "Install"} */}
      </>
      {/* )} */}
    </Button>
  );
};

export default ConnectWalletButton;
