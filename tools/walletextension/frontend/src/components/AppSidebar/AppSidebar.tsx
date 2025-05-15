import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from '@/components/ui/sidebar';
import { ChevronDown, Search, Settings } from 'lucide-react';
import { TbHexagons, TbPigMoney, TbUniverse } from 'react-icons/tb';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { RiPokerSpadesFill } from 'react-icons/ri';
import { GiCableStayedBridge } from 'react-icons/gi';
import { LuTextSearch } from 'react-icons/lu';
import ConnectWalletButton from '@/components/ConnectWallet/ConnectWalletButton';
import AccountSettings from '@/components/AccountSettings/AccountSettings';
import { GrDocumentText } from 'react-icons/gr';
import { CgWebsite } from 'react-icons/cg';

export function AppSidebar() {
    const zenItems = [
        {
            title: 'House of TEN',
            url: 'https://houseof.ten.xyz',
            icon: RiPokerSpadesFill,
        },
        {
            title: 'Battleships',
            url: 'https://battleships.ten.xyz',
            icon: TbHexagons,
        },
        {
            title: 'TENZEN',
            url: 'https://tenzen.ten.xyz',
            icon: TbUniverse,
        },
    ];
    const aItems = [
        {
            title: 'Bridge',
            url: 'https://bridge.ten.xyz',
            icon: GiCableStayedBridge,
        },
        {
            title: 'TenScan',
            url: 'https://scan.ten.xyz',
            icon: LuTextSearch,
        },
        {
            title: 'staking',
            url: 'https://staking.ten.xyz',
            icon: TbPigMoney,
        },
    ];
    const learnItems = [
        {
            title: 'Documentation',
            url: 'https://docs.ten.xyz',
            icon: GrDocumentText,
        },
        {
            title: 'Blog',
            url: 'https://blog.ten.xyz',
            icon: CgWebsite,
        },
    ];

    return (
        <Sidebar>
            <SidebarHeader className="mt-24">{/*<ConnectWalletButton />*/}</SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            {aItems.map((item) => (
                                <SidebarMenuItem key={item.title}>
                                    <SidebarMenuButton asChild>
                                        <a href={item.url}>
                                            <item.icon />
                                            <span>{item.title}</span>
                                        </a>
                                    </SidebarMenuButton>
                                </SidebarMenuItem>
                            ))}
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>

                <Collapsible defaultOpen className="group/collapsible">
                    <SidebarGroup>
                        <SidebarGroupLabel asChild>
                            <CollapsibleTrigger className="w-full">
                                ZEN Earners
                                <ChevronDown className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
                            </CollapsibleTrigger>
                        </SidebarGroupLabel>
                        <CollapsibleContent>
                            <SidebarGroupContent>
                                <SidebarMenu>
                                    {zenItems.map((item) => (
                                        <SidebarMenuItem key={item.title}>
                                            <SidebarMenuButton asChild>
                                                <a href={item.url}>
                                                    <item.icon />
                                                    <span>{item.title}</span>
                                                </a>
                                            </SidebarMenuButton>
                                        </SidebarMenuItem>
                                    ))}
                                </SidebarMenu>
                            </SidebarGroupContent>
                        </CollapsibleContent>
                    </SidebarGroup>
                </Collapsible>

                <Collapsible defaultOpen className="group/collapsible">
                    <SidebarGroup>
                        <SidebarGroupLabel asChild>
                            <CollapsibleTrigger className="w-full">
                                Learn
                                <ChevronDown className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
                            </CollapsibleTrigger>
                        </SidebarGroupLabel>
                        <CollapsibleContent>
                            <SidebarGroupContent>
                                <SidebarMenu>
                                    {learnItems.map((item) => (
                                        <SidebarMenuItem key={item.title}>
                                            <SidebarMenuButton asChild>
                                                <a href={item.url}>
                                                    <item.icon />
                                                    <span>{item.title}</span>
                                                </a>
                                            </SidebarMenuButton>
                                        </SidebarMenuItem>
                                    ))}
                                </SidebarMenu>
                            </SidebarGroupContent>
                        </CollapsibleContent>
                    </SidebarGroup>
                </Collapsible>

                <Collapsible className="group/collapsible">
                    <SidebarGroup>
                        <SidebarGroupLabel asChild>
                            <CollapsibleTrigger className="w-full">
                                Tools
                                <ChevronDown className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
                            </CollapsibleTrigger>
                        </SidebarGroupLabel>
                        <CollapsibleContent>
                            <SidebarGroupContent>
                                <SidebarMenu></SidebarMenu>
                            </SidebarGroupContent>
                        </CollapsibleContent>
                    </SidebarGroup>
                </Collapsible>

                <Collapsible className="group/collapsible">
                    <SidebarGroup>
                        <SidebarGroupLabel asChild>
                            <CollapsibleTrigger className="w-full">
                                Community
                                <ChevronDown className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
                            </CollapsibleTrigger>
                        </SidebarGroupLabel>
                        <CollapsibleContent>
                            <SidebarGroupContent>
                                <SidebarMenu></SidebarMenu>
                            </SidebarGroupContent>
                        </CollapsibleContent>
                    </SidebarGroup>
                </Collapsible>
            </SidebarContent>
            <SidebarFooter></SidebarFooter>
        </Sidebar>
    );
}
