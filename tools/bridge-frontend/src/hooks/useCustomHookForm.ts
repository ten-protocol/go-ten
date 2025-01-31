import { yupResolver } from "@hookform/resolvers/yup";
import { useForm } from "react-hook-form";

const useCustomHookForm = (formSchema: any, options?: any) => {
  return useForm({
    resolver: yupResolver(formSchema),
    mode: "onBlur",
    ...options,
  });
};

export default useCustomHookForm;
