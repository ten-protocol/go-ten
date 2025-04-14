import {useTenChainAuth} from "@/hooks/useTenChainAuth";
import {useEffect} from "react";
import {shortenAddress} from "@/lib/utils";
import classNames from 'classnames';
import {Badge} from "@/components/ui/badge";
import {Button} from "@/components/ui/button";
import {FaKey} from "react-icons/fa6";
import {Tooltip, TooltipContent, TooltipTrigger} from "@/components/ui/tooltip";
import {BarLoader} from "react-spinners";

type Props = {
    address: `0x${string}`
    active?: boolean
}

export default function ConnectedAccount({address, active}: Props) {
    const {isAuthenticated, isAuthenticatedLoading, authenticateAccount} = useTenChainAuth(address);

    useEffect(() => {
        if (!isAuthenticated && !isAuthenticatedLoading) {
            // authenticateAccount();
        }
    }, [isAuthenticatedLoading]);

    const statusClasses = classNames('rounded-full w-3 h-3 inline-block mr-2', {
        "bg-green-500": active && isAuthenticated,
        "bg-red-500": active && !isAuthenticated,
        "border border-green-500": !active && isAuthenticated,
        "border border-red-500": !active && !isAuthenticated,
    });

    const authClasses = classNames('', {
        "bg-green-500": isAuthenticated,
        "bg-red-500": !isAuthenticated,
    });


    return (

        <div className="bg-white dark:bg-neutral-800 p-3 rounded-md mb-3 flex gap-4 justify-between items-center">
            <p className="text-sm break-all">
                <span className={statusClasses}/>
                <span className="font-bold">{shortenAddress(address)}</span>
            </p>
            {isAuthenticatedLoading ? <BarLoader color="white" />
                :
            <div className="flex items-center gap-4">
                <Badge className={authClasses}>
                    {isAuthenticated ? "TRUE" : "FALSE"}
                </Badge>
                {!isAuthenticated &&
                    <Tooltip>
                    <TooltipTrigger asChild><Button size="icon" onClick={authenticateAccount} className="cursor-pointer">
                        <FaKey />
                    </Button></TooltipTrigger>
                    <TooltipContent>
                        Authenticate Account
                    </TooltipContent>
                </Tooltip>
                }
            </div>}
        </div>
    )
}