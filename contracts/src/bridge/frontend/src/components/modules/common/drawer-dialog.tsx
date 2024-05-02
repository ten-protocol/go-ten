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

export function DrawerDialog({
  FormComponent,
}: {
  FormComponent: React.JSXElementConstructor<any>;
}) {
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
            <small>Transfer to a different address</small>
          </Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Enter Transfer Address</DialogTitle>
            <DialogDescription>
              This address will be used to transfer the asset to.
            </DialogDescription>
          </DialogHeader>
          <FormComponent />
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
          <small>Transfer to a different address</small>
        </Button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Enter Transfer Address</DrawerTitle>
          <DrawerDescription>
            This address will be used to transfer the asset to.
          </DrawerDescription>
        </DrawerHeader>
        <FormComponent className="px-4" />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}
