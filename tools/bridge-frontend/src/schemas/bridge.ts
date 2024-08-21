import * as yup from "yup";

export const bridgeSchema = yup.object().shape({
  amount: yup.string().required("Amount is required"),
  fromChain: yup.string().required("From Chain is required"),
  toChain: yup.string().required("To Chain is required"),
  token: yup.string().required("Select a token"),
  receiver: yup.string().required("Receiver is required"),
});
