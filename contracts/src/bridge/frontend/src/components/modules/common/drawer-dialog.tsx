import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogTitle,
  DialogDescription,
} from "../../ui/dialog";
import React from "react";
import { Button } from "../../ui/button";
import { DialogHeader } from "../../ui/dialog";
import {
  Drawer,
  DrawerTrigger,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerDescription,
  DrawerFooter,
  DrawerClose,
} from "../../ui/drawer";
import { useMediaQuery } from "@/src/hooks/useMediaQuery";
import { PlusIcon } from "@radix-ui/react-icons";
import { Label } from "../../ui/label";
import { Input } from "../../ui/input";
import { Alert, AlertDescription } from "../../ui/alert";
import { Terminal } from "lucide-react";
import { cn } from "@/src/lib/utils";
import { useFormHook } from "@/src/hooks/useForm";

export function DrawerDialog() {
  const { form } = useFormHook();
  const receiver = form.watch("receiver");
  const [open, setOpen] = React.useState(false);
  const isDesktop = useMediaQuery("(min-width: 768px)");

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button
            className="text-sm font-bold leading-none hover:text-primary hover:bg-transparent"
            variant="ghost"
          >
            <PlusIcon className="h-3 w-3 mr-1" />
            <small>Edit destination address</small>
          </Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Enter Transfer Address</DialogTitle>
            <DialogDescription>
              This address will be used to transfer the asset to.
            </DialogDescription>
          </DialogHeader>
          <FormComponent setOpen={setOpen} form={form} receiver={receiver} />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerTrigger asChild>
        <Button
          className="text-sm font-bold leading-none hover:text-primary hover:bg-transparent"
          variant="ghost"
        >
          <PlusIcon className="h-3 w-3 mr-1" />
          <small>Edit destination address</small>
        </Button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Enter Transfer Address</DrawerTitle>
          <DrawerDescription>
            This address will be used to transfer the asset to.
          </DrawerDescription>
        </DrawerHeader>
        <FormComponent
          className="px-4"
          setOpen={setOpen}
          form={form}
          receiver={receiver}
        />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}

const FormComponent = ({ className, setOpen, form, receiver }: any) => {
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
};
