import type { Metadata } from 'next';
import { Chakra_Petch, Geist, Geist_Mono } from 'next/font/google';
import './globals.scss';
import { Providers } from '@/providers';

const geistSans = Geist({
    variable: '--font-geist-sans',
    subsets: ['latin'],
});

const geistMono = Geist_Mono({
    variable: '--font-geist-mono',
    subsets: ['latin'],
});

const chakra = Chakra_Petch({
    subsets: ['latin'],
    variable: '--font-chakra',
    display: 'swap',
    weight: ['300', '400', '500', '600', '700'],
});

export const metadata: Metadata = {
    title: 'Next.js RainbowKit Demo',
    description: 'A Next.js app with RainbowKit wallet integration',
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" className="dark" style={{colorScheme: "dark"}}>
            <body
                className={`${geistSans.variable} ${geistMono.variable} ${chakra.variable} antialiased`}
            >
                <Providers>{children}</Providers>
            </body>
        </html>
    );
}
