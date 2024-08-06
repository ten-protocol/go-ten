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

export const ChainSelect = ({
  form,
  chains,
  name,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  chains: Chain[];
  name: string;
}) => {
  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => (
        <FormItem>
          <Select defaultValue={field.value} onValueChange={field.onChange}>
            <FormControl>
              <SelectTrigger className="h-8 bg-muted">
                <span>{field.value || "Select Chain"}</span>
                {/* <SelectValue placeholder={field.value || "Select Chain"}  */}
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              {chains.map((chain: Chain) => (
                <SelectItem
                  key={chain.value}
                  value={chain.value}
                  disabled={!chain.isEnabled}
                >
                  {chain.value}
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
