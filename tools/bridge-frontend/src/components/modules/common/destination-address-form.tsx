import { cn } from "@/src/lib/utils";
import { Terminal } from "lucide-react";
import React from "react";
import { Alert, AlertDescription } from "../../ui/alert";
import { Button } from "../../ui/button";
import { Label } from "../../ui/label";
import { Input } from "../../ui/input";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";
import { handleStorage } from "@/src/lib/utils";

export default function FormComponent({
  form,
  className,
  setOpen,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  className?: string;
  setOpen: (value: boolean) => void;
}) {
  const receiver = form.getValues("receiver");
  const addAddressToMainForm = (e: any) => {
    e.preventDefault();
    const address = e.target.form[0].value;
    form.setValue("receiver", address);
    handleStorage.save("tenBridgeReceiver", address);
    setOpen(false);
  };
  return (
    <form className={cn("grid items-start gap-4", className)}>
      <div className="grid gap-2">
        <Label htmlFor="address">Address</Label>
        <Input
          type="address"
          id="address"
          placeholder="Enter address"
          defaultValue={receiver}
        />
      </div>
      <Alert variant={"warning"} className="flex items-center space-x-2">
        <Terminal className="h-4 w-4" />
        <AlertDescription>
          Make sure the address is correct before submitting.
        </AlertDescription>
      </Alert>
      <Button type="button" onClick={addAddressToMainForm}>
        Add destination address
      </Button>
    </form>
  );
}
