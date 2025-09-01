import {RecentRollups} from "@/src/components/modules/dashboard/recent-rollups";
import {pageLinks} from "@/src/routes";
import React from "react";
import {ResponseDataInterface} from "@repo/ui/lib/types/common";
import {cn} from "@repo/ui/lib/utils";
import {Skeleton} from "@repo/ui/components/shared/skeleton";
import {Button} from "@repo/ui/components/shared/button";
import Link from "next/link";

type Props = {
    title: string,
    data: ResponseDataInterface<any>,
    goTo: string,
    className: string,
    component: any
}

export default function RecentCard({title, goTo, className, data, component}: Props) {


   console.log(data)
    return (
        <div className="recents-card-shape h-[450px] relative col-span-1 md:col-span-2 lg:col-span-3 overflow-hidden">
            <div className="absolute inset-[1px] pointer-events-none">
                <div className="absolute inset-0 animate-scan-overlay"/>
            </div>
            <header className="flex justify-between p-4">
                <h3 className="opacity-80 text-lg">{title}</h3>
                <Link
                    href={{
                        pathname: goTo,
                    }}
                >
                    <Button variant="outline" size="sm">
                        View All
                    </Button>
                </Link>
            </header>
            <div className="h-[375px] overflow-y-scroll pl-2 mr-2 overflow-x-hidden" style={{
                maskImage: 'linear-gradient(to bottom, black 0%, black 80%, transparent 100%)',
                WebkitMaskImage: 'linear-gradient(to bottom, black 0%, black 80%, transparent 100%)'
            }}>
                {data ? (
                    component
                ) : (
                    <Skeleton className="w-full h-[200px] rounded-lg" />
                )}
            </div>
        </div>
    )
}