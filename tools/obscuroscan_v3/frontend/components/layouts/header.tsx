import Image from "next/image";
import { MainNav } from "../main-nav";
import { ModeToggle } from "../mode-toggle";
import ConnectWalletButton from "../modules/common/connect-wallet";
import { Search } from "../search";

export default function Header() {
  return (
    <div className="border-b">
      <div className="flex h-16 items-center px-4">
        <Image
          src="/assets/images/obscuro_black.png"
          width={200}
          height={32}
          alt="Obscuro Logo"
        />
        <MainNav className="mx-6" />
        <div className="ml-auto flex items-center space-x-4">
          <Search />
          <ModeToggle />
          <ConnectWalletButton />
        </div>
      </div>
    </div>
  );
}
