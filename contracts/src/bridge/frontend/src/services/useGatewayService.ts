import { ToastType } from "../types";
import { showToast } from "../components/ui/use-toast";
import { fetchTestnetStatus } from "@/api/general";

const useGatewayService = () => {
  const getTestnetStatus = async () => {
    try {
      return await fetchTestnetStatus();
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to connect to Ten Testnet");
      throw error;
    }
  };

  return {
    getTestnetStatus,
  };
};

export default useGatewayService;
