import { showToast } from "@/components/ui/use-toast";
import { isValidTokenFormat } from ".";
import { ToastType } from "@/types/interfaces";

export const validateToken = (token: string) => {
  if (!isValidTokenFormat(token)) {
    showToast(ToastType.INFO, "Invalid token format. Please refresh the page.");
    throw new Error("Invalid token format");
  }
};

export const handleError = (error: any, message: string) => {
  console.error(error);
  showToast(ToastType.DESTRUCTIVE, `${message}: ${error?.message || error}`);
};
