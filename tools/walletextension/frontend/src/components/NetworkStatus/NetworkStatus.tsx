import {fetchTestnetStatus} from "@/api/general";
import {useQuery} from "@tanstack/react-query";
import {Badge} from "@/components/ui/badge";
import classNames from "classnames";

export default function NetworkStatus() {
    const {data, isLoading, isError} = useQuery({queryKey: ['networkStatus'], queryFn: fetchTestnetStatus});
    const networkIsLive = !!data?.result?.OverallHealth
    const badgeClasses = classNames('', {
        'bg-red-500': !networkIsLive && !isLoading,
        'bg-green-400': networkIsLive,
    })

    return (
        <div className="flex gap-2 items-center">
            <p className="opacity-70 text-sm">Network Status:</p>
            <Badge className={badgeClasses} variant={isLoading ? 'outline' : 'default'}>
                {isLoading && 'Checking'}
                {networkIsLive ? 'LIVE' : 'DOWN'}
            </Badge>
        </div>
    )
}