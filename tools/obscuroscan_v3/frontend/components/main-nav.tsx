import React from "react";
import Link from "next/link";

import { cn } from "@/src/lib/utils";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
} from "./ui/dropdown-menu";

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <nav
      className={cn("flex items-center space-x-4 lg:space-x-6", className)}
      {...props}
    >
      <Link
        href="/"
        className="text-sm font-medium transition-colors hover:text-primary"
      >
        Home
      </Link>
      <Link
        href="/personal"
        className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
      >
        Personal
      </Link>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="clear"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            Blockchain
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56" align="end" forceMount>
          <DropdownMenuGroup>
            <Link href="/transactions">
              <DropdownMenuItem>Transactions</DropdownMenuItem>
            </Link>
            <Link href="/blocks">
              <DropdownMenuItem>Blocks</DropdownMenuItem>
            </Link>
            <Link href="/batches">
              <DropdownMenuItem>Batches</DropdownMenuItem>
            </Link>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="clear"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            Resources
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56" align="end" forceMount>
          <DropdownMenuGroup>
            <Link href="/resources/decrypt">
              <DropdownMenuItem>Decrypt</DropdownMenuItem>
            </Link>
            <Link href="/resources/verify">
              <DropdownMenuItem>Verified Data</DropdownMenuItem>
            </Link>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </nav>
  );
}
