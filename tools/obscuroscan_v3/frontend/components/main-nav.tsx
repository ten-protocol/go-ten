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

import { ChevronDownIcon } from "@radix-ui/react-icons";
import { NavLinks } from "@/src/routes";
import { NavLink } from "@/src/types/interfaces";

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <nav
      className={cn("flex items-center space-x-4 lg:space-x-6", className)}
      {...props}
    >
      {NavLinks.map((navLink: NavLink, index: number) => {
        return navLink.isDropdown ? (
          <DropdownMenu key={index}>
            <DropdownMenuTrigger asChild>
              <Button
                variant="clear"
                className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary p-0"
              >
                {navLink.label} <ChevronDownIcon className="ml-1 h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56" align="end" forceMount>
              <DropdownMenuGroup>
                {navLink.subNavLinks &&
                  navLink.subNavLinks.map((subNavLink: NavLink) =>
                    subNavLink.isExternal ? (
                      <a href={subNavLink.href} key={subNavLink.label}>
                        <DropdownMenuItem>{subNavLink.label}</DropdownMenuItem>
                      </a>
                    ) : (
                      subNavLink.href && (
                        <Link href={subNavLink.href} key={subNavLink.label}>
                          <DropdownMenuItem className="cursor-pointer">
                            {subNavLink.label}
                          </DropdownMenuItem>
                        </Link>
                      )
                    )
                  )}
              </DropdownMenuGroup>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : navLink.isExternal ? (
          <a
            key={index}
            href={navLink.href}
            className="text-sm font-medium transition-colors hover:text-primary"
          >
            {navLink.label}
          </a>
        ) : (
          navLink.href && (
            <Link
              key={index}
              href={navLink.href}
              className="text-sm font-medium transition-colors hover:text-primary"
            >
              {navLink.label}
            </Link>
          )
        );
      })}
    </nav>
  );
}
