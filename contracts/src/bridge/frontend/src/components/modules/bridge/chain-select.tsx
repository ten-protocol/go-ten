import { Chain } from "@/src/types";
import { FormField, FormItem, FormControl, FormMessage } from "../../ui/form";
import {
  Select,
  SelectTrigger,
  SelectContent,
  SelectItem,
  SelectValue,
} from "../../ui/select";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";

export const ChainSelectFrom = ({
  form,
  chains,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  chains: Chain[];
}) => {
  return (
    <FormField
      control={form.control}
      name="fromChain"
      render={({ field }) => (
        <FormItem>
          <Select defaultValue={field.value} onValueChange={field.onChange}>
            <FormControl>
              <SelectTrigger className="h-8 bg-muted">
                <SelectValue placeholder={field.value || "Select Chain"} />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              {chains.map((chain: Chain) => (
                <SelectItem
                  key={chain.value}
                  value={chain.value}
                  disabled={!chain.isEnabled}
                >
                  {chain.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};

export const ChainSelectTo = ({
  form,
  chains,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  chains: Chain[];
}) => {
  return (
    <FormField
      control={form.control}
      name="toChain"
      render={({ field }) => (
        <FormItem>
          <Select defaultValue={field.value} onValueChange={field.onChange}>
            <FormControl>
              <SelectTrigger className="h-8 bg-muted">
                <SelectValue placeholder={field.value || "Select Chain"} />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              {chains.map((chain: Chain) => (
                <SelectItem
                  key={chain.value}
                  value={chain.value}
                  disabled={!chain.isEnabled}
                >
                  {chain.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
