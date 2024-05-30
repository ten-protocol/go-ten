import { cn } from "@/src/lib/utils";
import { Terminal } from "lucide-react";
import React from "react";
import { Alert, AlertDescription } from "../../ui/alert";
import { Button } from "../../ui/button";
import { Label } from "../../ui/label";
import { Input } from "../../ui/input";

export default function FormComponent({
  className,
  setOpen,
  form,
  receiver,
}: {
  className?: string;
  setOpen: (value: boolean) => void;
  form: any;
  receiver: string | undefined;
}) {
  const addAddressToMainForm = (e: any) => {
    e.preventDefault();
    form.setValue("receiver", e.target.elements.address.value);
    setOpen(false);
  };
  return (
    <form
      className={cn("grid items-start gap-4", className)}
      onSubmit={addAddressToMainForm}
    >
      <div className="grid gap-2">
        <Label htmlFor="address">Address</Label>
        <Input
          type="address"
          id="address"
          placeholder={receiver || "Enter address"}
          defaultValue={receiver}
        />
      </div>
      <Alert variant={"warning"} className="flex items-center space-x-2">
        <Terminal className="h-4 w-4" />
        <AlertDescription>
          Make sure the address is correct before submitting.
        </AlertDescription>
      </Alert>
      <Button type="submit">Add destination address</Button>
    </form>
  );
}
