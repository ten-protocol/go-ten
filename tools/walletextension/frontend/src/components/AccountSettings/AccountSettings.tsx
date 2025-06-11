import { IoMdSettings } from 'react-icons/io';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import RevokeAccount from '@/components/AccountSettings/RevokeAccount';
import ViewTenToken from '@/components/AccountSettings/ViewTenToken';
import { useState } from 'react';
import { useDisconnect } from 'wagmi';

export default function AccountSettings() {
    const { disconnect } = useDisconnect();
    const [isRevokeOpen, setIsRevokeOpen] = useState(false);
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [isViewPrivateKeyOpen, setIsViewPrivateKeyOpen] = useState(false);

    const handleOpenRevoke = () => {
        setIsDropdownOpen(false);
        setIsRevokeOpen(true);
    };

    const handleOpenViewPrivateKey = () => {
        setIsDropdownOpen(false);
        setIsViewPrivateKeyOpen(true);
    };

    const handleDisconnect = () => {
        disconnect();
    };

    return (
        <div>
            <RevokeAccount isOpen={isRevokeOpen} onChange={setIsRevokeOpen} />
            <ViewTenToken isOpen={isViewPrivateKeyOpen} onChange={setIsViewPrivateKeyOpen} />
            <DropdownMenu open={isDropdownOpen} onOpenChange={setIsDropdownOpen}>
                <DropdownMenuTrigger asChild>
                    <Button variant="outline">
                        <IoMdSettings />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56">
                    <DropdownMenuLabel>My Account</DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuGroup>
                        <DropdownMenuItem onClick={handleOpenViewPrivateKey}>
                            View TEN Token
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={handleDisconnect}>
                            Disconnect Wallet
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={handleOpenRevoke}>
                            Revoke Account
                        </DropdownMenuItem>
                    </DropdownMenuGroup>
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    );
}
