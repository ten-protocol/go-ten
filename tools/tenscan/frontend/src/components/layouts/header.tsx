import { MainNav } from "../main-nav";
import { ModeToggle } from "../mode-toggle";
import ConnectWalletButton from "../modules/common/connect-wallet";
import Link from "next/link";
import { HamburgerMenuIcon } from "@radix-ui/react-icons";
import { X } from "lucide-react";
import { useState } from "react";
import { Button } from "../ui/button";
import HealthIndicator from "../health-indicator";
import Image from "next/image";

export default function Header() {
  return (
    <div className="border-b">
      <div className="flex h-16 justify-between items-center px-4">
        <Link href="/" className="min-w-[80px]">
          <Image
            src="/assets/images/black_logotype.png"
            alt="Logo"
            width={150}
            height={40}
            className="cursor-pointer dark:hidden"
          />
          <Image
            src="/assets/images/white_logotype.png"
            alt="Logo"
            width={150}
            height={40}
            className="cursor-pointer hidden dark:block"
          />
        </Link>
        <div className="hidden md:flex items-center space-x-2">
          <MainNav className="px-2" />
          <div className="flex items-center space-x-4">
            <HealthIndicator />
            <ModeToggle />
            <ConnectWalletButton />
          </div>
        </div>
        <div className="flex items-center space-x-4 md:hidden">
          <MobileMenu />
        </div>
      </div>
    </div>
  );
}

const MobileMenu = () => {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className="relative flex items-center space-x-">
      <HealthIndicator />
      <Button
        variant={"clear"}
        className="text-muted-foreground hover:text-primary transition-colors"
        onClick={() => setMenuOpen(!menuOpen)}
      >
        {menuOpen ? <X /> : <HamburgerMenuIcon />}
      </Button>
      {menuOpen && (
        <div className="absolute z-10 top-0 right-0 mt-12">
          <div className="bg-background border rounded-lg shadow-lg">
            <div className="flex flex-col p-4 space-y-2">
              <MainNav className="flex flex-col" />
              <ConnectWalletButton />
              <ModeToggle />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
