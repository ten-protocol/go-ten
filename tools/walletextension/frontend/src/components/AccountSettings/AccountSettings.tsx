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
import { useState } from 'react';

export default function AccountSettings() {
    const [isRevokeOpen, setIsRevokeOpen] = useState(false);
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);

    const handleOpenRevoke = () => {
        setIsDropdownOpen(false);
        setIsRevokeOpen(true);
    };

    return (
        <div>
            <RevokeAccount isOpen={isRevokeOpen} onChange={setIsRevokeOpen} />
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
                        <DropdownMenuItem>View Private Key</DropdownMenuItem>
                        <DropdownMenuItem>Disconnect Wallet</DropdownMenuItem>
                        <DropdownMenuItem onClick={handleOpenRevoke}>
                            Revoke Account
                        </DropdownMenuItem>
                    </DropdownMenuGroup>
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    );
}
