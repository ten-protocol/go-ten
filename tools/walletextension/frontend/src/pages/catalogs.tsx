"use client";
import Image from "next/image";
import React, { ReactNode } from "react";
import { BookmarkIcon } from "lucide-react";
import { Button } from "@repo/ui/components/shared/button";

const trends = [
  {
    cover: { url: "", alt: "" },
    thumbnail: { url: "", alt: "" },
    title: "Dookey Dash: Unclogged",
    description:
      "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Ab earum in nemo recusandae tempore! At distinctio eius ipsa ipsam mollitia nemo nostrum quis quo reprehenderit vel. Deserunt et nam rerum.",
  },
  {
    cover: { url: "", alt: "" },
    thumbnail: { url: "", alt: "" },
    title: "GU Factory",
    description:
      "By trench warriors, for trech warriors. Arbitrum One and Orbit token factory built on top of",
  },
];

function Trends() {
  return (
    <section>
      <Heading
        heading={"Trending Projects"}
        subHeading={"The most popular dApps on TEN"}
      />

      <div
        className={
          "grid mt-4 gap-4 grid-[repeat(auto-fill,_minmax(200px,_400px))]"
        }
        style={{
          gridTemplateColumns: "repeat(auto-fill, minmax(150px, 400px))",
        }}
      >
        {trends.map((item) => {
          return <BaseCard data={item} key={item.title} />;
        })}
      </div>
    </section>
  );
}

function Ecosystem() {
  const buttons = [
    { text: "DeFi" },
    { text: "Bridges" },
    { text: "Gaming" },
    { text: "NFTs" },
    { text: "Infra & Tools" },
  ];
  const active = buttons[0];

  return (
    <section>
      <Heading heading={"Ecosystem Essentials"}>
        <div className={"flex gap-x-2 mt-4 mb-2"}>
          {buttons.map((e) => {
            return (
              <Button
                key={e.text}
                variant={active === e ? "default" : "secondary"}
                size={"sm"}
              >
                {e.text}
              </Button>
            );
          })}
        </div>
      </Heading>

      <div
        className={
          "grid mt-4 gap-4 grid-[repeat(auto-fill,_minmax(200px,_400px))]"
        }
        style={{
          gridTemplateColumns: "repeat(auto-fill, minmax(150px, 400px))",
        }}
      >
        {trends.map((item) => {
          return <BaseCard data={item} key={item.title} />;
        })}
      </div>
    </section>
  );
}

function BaseCard({ data: item }: { data: (typeof trends)[0] }) {
  return (
    <div
      className={
        "aspect-[333/184] flex flex-col overflow-hidden justify-end rounded-xl relative bg-[#111111]"
      }
    >
      <Image {...item.cover} fill className={"object-cover object-center"} />
      <div
        className={
          "inset-0 z-20 absolute bg-gradient-to-t from-white/[0.4] to-transparent"
        }
      />
      <div
        className={
          "p-4 flex-1 gap-2 flex flex-col gap-2 border border-black z-20"
        }
      >
        <div className={"flex text-sm justify-end -mt-2 -mr-2"}>
          <button
            className={
              "inline-flex hover:bg-white/[0.2] backdrop-blur filter rounded-full items-center justify-center w-10 h-10"
            }
          >
            <BookmarkIcon size={"1.2em"} className={"opacity-80"} />
          </button>
        </div>
        <div />
        <Image
          {...item.thumbnail}
          className={
            "bg-white/[0.8] w-12 rounded-lg shadow-lg aspect-square backdrop-blur filter"
          }
        />
        <h2 className={"text-xl"}>
          <span className={"leading-[1ex]"}>{item.title}</span>
        </h2>

        <div className={"text-sm"}>
          <p
            className={
              "overflow-hidden overflow-ellipsis h-[4.8ex] leading-[2.4ex]"
            }
          >
            {item.description}
          </p>
        </div>
      </div>
    </div>
  );
}

function Heading({
  heading,
  subHeading,
  children,
}: {
  heading: string;
  subHeading?: string;
  children?: ReactNode;
}) {
  return (
    <>
      <hgroup>
        <h2 className={"text-2xl"}>{heading}</h2>
        {subHeading ? (
          <p className={"text-muted-foreground"}>{subHeading}</p>
        ) : null}
        {children}
      </hgroup>
      <div className={"border-b pb-2"} />
    </>
  );
}

export function MarketPlace() {
  return (
    <div className={"flex flex-col gap-12"}>
      <Trends />
      <Ecosystem />
    </div>
  );
}
