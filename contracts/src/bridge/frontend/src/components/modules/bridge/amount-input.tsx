import { FormField, FormItem, FormControl, FormMessage } from "../../ui/form";
import { Input } from "../../ui/input";

export const AmountInput = ({
  form,
  walletConnected,
}: {
  form: any;
  walletConnected: boolean;
}) => {
  return (
    <FormField
      control={form.control}
      name="amount"
      render={({ field }) => (
        <FormItem>
          <FormControl>
            <Input
              type="number"
              placeholder="0"
              className="text-2xl font-bold w-full dark:bg-[#292929] overflow-ellipsis"
              disabled={!walletConnected}
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
