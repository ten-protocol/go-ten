import { fetchTestnetStatus } from '@/api/general';
import { useQuery } from '@tanstack/react-query';
import { Badge } from '@/components/ui/badge';
import classNames from 'classnames';
import { tenNetworkName } from '@/lib/constants';
import { Loader2 } from 'lucide-react';

export default function NetworkStatus() {
    const { data, isLoading } = useQuery({
        queryKey: ['networkStatus'],
        queryFn: fetchTestnetStatus,
    });
    const networkIsLive = !!data?.result?.OverallHealth;
    const badgeClasses = classNames('', {
        'bg-red-500': !networkIsLive && !isLoading,
        'bg-green-400': networkIsLive,
    });

    return (
        <div className="flex gap-2 items-center">
            <p className="hidden md:block opacity-70 text-sm whitespace-nowrap">
                {tenNetworkName} Status:
            </p>
            <Badge className={badgeClasses} variant={isLoading ? 'outline' : 'default'}>
                {isLoading && <Loader2 className="h-20 w-20 animate-spin" />}
                {!isLoading ? (networkIsLive ? 'LIVE' : 'DOWN') : null}
            </Badge>
        </div>
    );
}
