import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useWalletStore } from "../components/providers/wallet-provider";
import { L1TOKENS, L2TOKENS } from "../lib/constants";

export const useFormHook = () => {
  const { fromChains, toChains, isL1ToL2 } = useWalletStore();

  const FormSchema = z.object({
    amount: z.string().nonempty({
      message: "Amount is required.",
    }),
    fromChain: z.string().nonempty({
      message: "From Chain is required.",
    }),
    toChain: z.string().nonempty({
      message: "To Chain is required.",
    }),
    token: z.string().nonempty({
      message: "Select a token.",
    }),
  });

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      amount: "",
      fromChain: fromChains[0].value,
      toChain: toChains[0].value,
      token: isL1ToL2 ? L1TOKENS[0].value : L2TOKENS[0].value,
    },
  });

  return { form, FormSchema };
};
