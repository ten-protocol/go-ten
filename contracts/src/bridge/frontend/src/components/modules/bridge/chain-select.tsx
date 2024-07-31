import { Chain } from "@/src/types";
import { FormField, FormItem, FormControl, FormMessage } from "../../ui/form";
import {
  Select,
  SelectTrigger,
  SelectContent,
  SelectItem,
  SelectValue,
} from "../../ui/select";

export const ChainSelect = ({
  form,
  chains,
  name,
}: {
  form: any;
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
                <SelectValue placeholder={field.value || "Select Chain"} />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              {chains.map((chain: any) => (
                <SelectItem key={chain.value} value={chain.value}>
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
