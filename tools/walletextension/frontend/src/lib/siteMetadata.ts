import type { Metadata } from 'next';

export const siteMetadata: Metadata = {
    title: 'TEN Gateway - Your Portal to the TEN Network',
    description:
        'Your secure entry point to the TEN ecosystem. Easily authenticate your accounts and access dApps built on TEN. Experience seamless, confidential blockchain interactions.',
    keywords:
        'TEN Network, TEN Gateway, privacy blockchain, Layer 2, confidential smart contracts, encrypted transactions, decentralized applications, dApps,secure DeFi, TEN ecosystem, wallet extension, crypto privacy',
    metadataBase: new URL('https://gateway.ten.xyz'),
    openGraph: {
        title: 'TEN Gateway - Your Portal to the TEN Network',
        description:
            'Your secure entry point to the TEN ecosystem. Easily authenticate your accounts and access dApps built on TEN. Experience seamless, confidential blockchain interactions.',
        url: 'https://gateway.ten.xyz',
        siteName: 'TEN Gateway',
        type: 'website',
        images: [
            {
                url: '/assets/og.png',
                width: 1200,
                height: 630,
                alt: 'TEN Gateway - Privacy-focused blockchain portal',
            },
        ],
    },
    twitter: {
        card: 'summary_large_image',
        site: '@tenprotocol',
        creator: '@tenprotocol',
        title: 'TEN Gateway - Your Portal to the TEN Network',
        description: 'Your secure entry point to the TEN ecosystem.',
        images: ['/assets/x.png'],
    },
    robots: {
        index: true,
        follow: true,
        googleBot: {
            index: true,
            follow: true,
            'max-video-preview': -1,
            'max-image-preview': 'large',
            'max-snippet': -1,
        },
    },
};
