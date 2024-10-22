import Image from "next/image";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarTrigger,
  useSidebar,
} from "@/components/ui/sidebar";
import {
  ChevronDown,
  Home,
  NetworkIcon,
  SendHorizonalIcon,
  Share2Icon,
  TrendingUpIcon,
} from "lucide-react";
import Link from "next/link";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { cn, safeArray } from "@repo/ui/lib/utils";

const items = [
  {
    title: "Home",
    url: "#",
    icon: Home,
  },
  {
    title: "Projects",
    url: "#",
    icon: Share2Icon,
    children: [
      { title: "DeFi", url: "/projects/defi" },
      { title: "Bridges and On-ramps", url: "/projects/onramps" },
      { title: "Gaming", url: "/projects/gaming" },
      { title: "NFTs", url: "/projects/nfts" },
      { title: "Infra & Tools", url: "/projects/infra-tools" },
    ],
  },
  {
    title: "Bridge",
    url: "#",
    icon: SendHorizonalIcon,
  },
  {
    title: "TenScan",
    url: "#",
    icon: SendHorizonalIcon,
  },
  {
    title: "Staking",
    url: "#",
    icon: SendHorizonalIcon,
  },
  {
    title: "ZEN",
    url: "#",
    icon: SendHorizonalIcon,
  },
  {
    title: "Learn",
    url: "#",
    icon: SendHorizonalIcon,
  },
  {
    title: "Tools",
    url: "#",
    icon: TrendingUpIcon,
  },
  {
    title: "Community",
    url: "#",
    icon: NetworkIcon,
  },
];

export function AppSidebar() {
  return (
    <Sidebar collapsible={"icon"}>
      <SidebarTrigger
        className={
          "border absolute transform right-0 translate-x-1/2 top-0 rounded-full translate-y-8 z-20 bg-[#1C1C1C]"
        }
      />

      <SidebarHeader>
        <Link href="/" className={"my-4"}>
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
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup className={"flex-1"}>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => {
                const subLinks = safeArray(item.children);

                return (
                  <Collapsible
                    key={item.title}
                    defaultOpen
                    className={"group/collapsible"}
                  >
                    <SidebarMenuItem>
                      <CollapsibleTrigger asChild>
                        <SidebarMenuButton asChild>
                          <a href={item.url}>
                            <item.icon />
                            <span>{item.title}</span>
                          </a>
                        </SidebarMenuButton>
                      </CollapsibleTrigger>

                      {subLinks.length > 0 ? (
                        <>
                          <CollapsibleTrigger asChild>
                            <SidebarMenuAction
                              className={
                                "data-[state=closed]:-rotate-90 transform"
                              }
                            >
                              <ChevronDown />
                            </SidebarMenuAction>
                          </CollapsibleTrigger>

                          <CollapsibleContent>
                            <SidebarMenuSub>
                              {subLinks.map((e: { title: string }) => {
                                return (
                                  <SidebarMenuSubItem key={e.title}>
                                    <SidebarMenuSubButton>
                                      {e.title}
                                    </SidebarMenuSubButton>
                                  </SidebarMenuSubItem>
                                );
                              })}
                            </SidebarMenuSub>
                          </CollapsibleContent>
                        </>
                      ) : null}
                    </SidebarMenuItem>
                  </Collapsible>
                );
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <AppSidebarFooter />
      </SidebarContent>
    </Sidebar>
  );
}

function AppSidebarFooter() {
  const { open } = useSidebar();

  return (
    <SidebarFooter>
      <div
        className={cn("pb-4 transform duration-200", {
          "opacity-0 translate-y-4 duration-[50ms]": !open,
          "opacity-100 translate-y-0 delay-200": open,
        })}
      >
        <div className={"text-base"}>
          The Final Missing Piece <br /> of Web3
        </div>
        <div className={"border-t my-1"} />
        <div className={"gap flex gap-x-2 text-xs text-gray-500"}>
          <span>ToS</span> • <span>Privacy Policy</span> •{" "}
          <span>Media Kit</span>
        </div>
      </div>
    </SidebarFooter>
  );
}
