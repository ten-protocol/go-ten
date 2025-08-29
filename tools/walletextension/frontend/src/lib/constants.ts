import {RiPokerSpadesFill} from "react-icons/ri";
import {TbHexagons, TbUniverse} from "react-icons/tb";
import {IoRocketSharp} from "react-icons/io5";
import {GiKoala} from "react-icons/gi";

export const tenGatewayAddress = process.env.NEXT_PUBLIC_GATEWAY_URL;

export const tenNetworkName = process.env.NEXT_PUBLIC_NETWORK_NAME || 'TEN Testnet';

export const tenscanAddress = process.env.NEXT_PUBLIC_TENSCAN_URL || 'https://tenscan.io';

export const socialLinks = {
    faucet: 'https://faucet.ten.xyz',
    github: 'https://github.com/ten-protocol',
    discord: 'https://discord.gg/tenprotocol',
    twitter: 'https://twitter.com/tenprotocol',
    twitterHandle: '@tenprotocol',
};

export const GOOGLE_ANALYTICS_ID = process.env.NEXT_PUBLIC_GOOGLE_ANALYTICS_ID;

export const SWITCHED_CODE = 4902;
export const tokenHexLength = 42;

export const tenGatewayVersion = 'v1';
export const tenChainIDDecimal = parseInt(process.env.NEXT_PUBLIC_CHAIN_ID || '443', 10);

export const environment = process.env.NEXT_PUBLIC_ENVIRONMENT;

export const tenChainIDHex = '0x' + tenChainIDDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'

export const METAMASK_CONNECTION_TIMEOUT = 3000;

export const userStorageAddress = '0x0000000000000000000000000000000000000001';

export const nativeCurrency = {
    name: 'Sepolia Ether',
    symbol: 'ETH',
    decimals: 18,
};

export const CONNECTION_STEPS = [
    'Hit Connect to TEN and start your journey',
    'Allow MetaMask to switch networks to the TEN Testnet',
    'Sign the <b>Signature Request</b> (this is not a transaction)',
];


export const PROMO_APPS = [
    {
        title: "Battleships",
        description: "Sink ships, win ZEN! Play a new vartiation of Battleships.",
        imageUrl: "/assets/promo/bs.png",
        url: "https://battleships.ten.xyz",
        icon: RiPokerSpadesFill,
    },
    {
        title: "House of TEN",
        description: "An Onchain poker tournament played by frontier AI models.",
        imageUrl: "/assets/promo/houseOfTen.png",
        url: "https://houseof.ten.xyz",
        icon: TbHexagons,
    },
    {
        title: "TENZEN",
        description: "Play to hit zero! Every extra zero boosts your prize!",
        imageUrl: "/assets/promo/tenzen.png",
        url: "https://tenzen.ten.xyz",
        icon: TbUniverse,
    },
    {
        title: "TEN X",
        description: "On chain crash games! Classic crash and PVP!",
        imageUrl: "/assets/promo/tenx.png",
        url: "https://tenx.ten.xyz",
        icon: IoRocketSharp,
    },
    {
        title: "Koala Counter",
        description: "TEN Mission to stop Wen TGE!",
        imageUrl: "/assets/promo/wenten.png",
        url: "https://wen.ten.xyz",
        icon: GiKoala,
    },
]
