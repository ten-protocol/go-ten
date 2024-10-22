import { MainNav } from "../main-nav";
import { ModeToggle } from "../mode-toggle";
import { HamburgerMenuIcon } from "@repo/ui/components/shared/react-icons";
import { useState } from "react";
import { Button } from "@repo/ui/components/shared/button";
import HealthIndicator from "../health-indicator";
import GatewayConnectButton from "../modules/connect-wallet/gateway-connect-button";
import { Input } from "@repo/ui/components/shared/input";
import { SearchIcon } from "lucide-react";

export default function Header() {
  return (
    <header>
      <div className="flex h-16 justify-between items-center px-4">
        <SearchInput />

        <div className="hidden md:flex items-center space-x-4">
          <MainNav className="mx-6" />
          <div className="flex items-center space-x-4">
            <HealthIndicator />
            <ModeToggle />
            <GatewayConnectButton />
          </div>
        </div>
        <div className="flex items-center space-x-4 md:hidden">
          <MobileMenu />
        </div>
      </div>
    </header>
  );
}

const MobileMenu = () => {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className="relative">
      <ModeToggle />
      <Button
        variant={"clear"}
        className="text-muted-foreground hover:text-primary transition-colors"
        onClick={() => setMenuOpen(!menuOpen)}
      >
        <HamburgerMenuIcon />
      </Button>
      {menuOpen && (
        <div className="absolute z-10 top-0 right-0 mt-12">
          <div className="bg-background border rounded-lg shadow-lg">
            <div className="flex flex-col p-4 space-y-2">
              <MainNav className="flex flex-col" />
              <GatewayConnectButton />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

function SearchInput() {
  return (
    <div className={"relative flex items-center text-xs"}>
      <SearchIcon className={"absolute mx-4"} size={"1em"} />
      <Input
        placeholder={"Search projects and more"}
        className={
          "indent-[2em] bg-[#1b1c1e] border-[#3c3d3f] min-w-[300px] text-[1em]"
        }
      />
    </div>
  );
}
