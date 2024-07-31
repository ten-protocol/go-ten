import { FormField, FormItem, FormControl, FormMessage } from "../../ui/form";
import {
  Select,
  SelectTrigger,
  SelectContent,
  SelectItem,
  SelectValue,
} from "../../ui/select";

export const TokenSelect = ({ form, tokens }: any) => {
  return (
    <FormField
      control={form.control}
      name="token"
      render={({ field }) => (
        <FormItem>
          <Select defaultValue={field.value} onValueChange={field.onChange}>
            <FormControl>
              <SelectTrigger className="h-8 dark:bg-[#292929]">
                <SelectValue placeholder={field.value || "Select Token"} />
              </SelectTrigger>
            </FormControl>
            <SelectContent side="top">
              {tokens.map((token: any) => (
                <SelectItem
                  key={token.value}
                  value={token.value}
                  disabled={!token.isEnabled}
                >
                  {token.value}
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
