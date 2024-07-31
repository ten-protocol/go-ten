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
import FormComponent from "./destination-address-form";

export function DrawerDialog({
  open,
  setOpen,
}: {
  open: boolean;
  setOpen: (value: boolean) => void;
}) {
  const isDesktop = useMediaQuery("(min-width: 768px)");

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Enter Transfer Address</DialogTitle>
            <DialogDescription>
              This address will be used to transfer the asset to.
            </DialogDescription>
          </DialogHeader>
          <FormComponent setOpen={setOpen} />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Enter Transfer Address</DrawerTitle>
          <DrawerDescription>
            This address will be used to transfer the asset to.
          </DrawerDescription>
        </DrawerHeader>
        <FormComponent className="px-4" setOpen={setOpen} />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}
