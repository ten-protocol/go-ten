import { ethers } from "ethers";
import { ethereum } from ".";
import { showToast } from "@/src/components/ui/use-toast";
import { IErrorMessages, ToastType } from "@/src/types";

export const getEthereumProvider = async () => {
  if (!ethereum) {
    throw new Error("No Ethereum provider detected");
  }
  const provider = new ethers.providers.Web3Provider(ethereum);
  return provider;
};

const errorMessages: Record<IErrorMessages, string> = {
  [IErrorMessages.UnknownAccount]:
    "Please ensure your wallet is unlocked and connected to the correct network",
  [IErrorMessages.InsufficientFunds]:
    "Insufficient funds. Please ensure you have enough balance to proceed",
  [IErrorMessages.UserDeniedTransactionSignature]:
    "Transaction rejected. Please sign the transaction to proceed",
  [IErrorMessages.UserRejectedTheRequest]:
    "Request rejected. Please try again with the correct permissions",
  [IErrorMessages.ExecutionReverted]:
    "Transaction reverted. Please check the transaction details and try again",
  [IErrorMessages.RateLimitExceeded]:
    "Rate limit exceeded. Please try again later",
  [IErrorMessages.WithdrawalSpent]:
    "Withdrawal already spent. Please check the transaction details",
};

export const handleError = (error: any, message: string) => {
  console.error(`Error: ${message}`, error);

  const errorReason = error?.reason || error?.message;

  if (errorReason in errorMessages) {
    showToast(
      ToastType.DESTRUCTIVE,
      errorMessages[errorReason as IErrorMessages]
    );
  } else {
    showToast(ToastType.DESTRUCTIVE, message);
  }

  throw new Error(message);
};
