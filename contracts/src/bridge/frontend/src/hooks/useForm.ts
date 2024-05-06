import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

export const useFormHook = () => {
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
    fromToken: z.string().nonempty({
      message: "From Token is required.",
    }),
    toToken: z.string().nonempty({
      message: "To Token is required.",
    }),
  });

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      amount: "",
      fromChain: "",
      toChain: "",
      fromToken: "",
      toToken: "",
    },
  });

  return { form, FormSchema };
};
