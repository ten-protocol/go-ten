# Next.js with RainbowKit

This is a [Next.js](https://nextjs.org/) project bootstrapped with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app) integrated with [RainbowKit](https://www.rainbowkit.com/) for easy Ethereum wallet connections.

## Prerequisites

You need to have the following installed:

- [Node.js](https://nodejs.org/) (v18 or higher)
- npm or yarn

## Getting Started

First, clone the repository and install the dependencies:

```bash
npm install
# or
yarn install
```

### Configuration

Before running the project, you need to:

1. Get a WalletConnect Project ID from [WalletConnect Cloud](https://cloud.walletconnect.com/)
2. Open `app/providers.tsx` and replace `YOUR_WALLETCONNECT_PROJECT_ID` with your actual project ID
3. Configure your custom chain in `app/providers.tsx` by updating the `customChain` object with your chain details

#### Custom Chain Configuration

The project includes a sample custom chain configuration. To connect to your own chain, update the following properties in the `customChain` object in `app/providers.tsx`:

```typescript
const customChain = {
  id: 1337, // Replace with your chain ID
  name: 'My Custom Chain', // Replace with your chain name
  nativeCurrency: {
    name: 'Custom Token', // Replace with your token name
    symbol: 'CTK', // Replace with your token symbol
    decimals: 18,
  },
  rpcUrls: {
    default: {
      http: ['https://your-custom-rpc-url.com'], // Replace with your RPC URL
    },
    public: {
      http: ['https://your-custom-rpc-url.com'], // Replace with your RPC URL
    },
  },
  blockExplorers: {
    default: {
      name: 'CustomScan', // Replace with your explorer name
      url: 'https://your-custom-explorer.com', // Replace with your block explorer URL
    },
  },
};
```

Common chain IDs:
- Local development chains: 1337 (Ganache), 31337 (Hardhat)
- Public testnets: 5 (Goerli), 80001 (Mumbai), 97 (BSC Testnet)
- Public mainnets: 1 (Ethereum), 137 (Polygon), 56 (Binance Smart Chain)

### Running the Development Server

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Features

- Next.js 14 App Router
- RainbowKit for easy wallet connections
- Tailwind CSS for styling
- TypeScript for type safety
- Support for multiple Ethereum networks (mainnet, sepolia, and custom chains)
- Chain switching UI

## Project Structure

- `app/components/ConnectButton.tsx` - Custom styled RainbowKit ConnectButton
- `app/components/ChainInfo.tsx` - UI for displaying and switching chains
- `app/providers.tsx` - RainbowKit and wagmi providers with custom chain config
- `app/page.tsx` - Main application page
- `app/layout.tsx` - Root layout with providers

## Customization

### Adding More Networks

To add more networks, edit `app/providers.tsx` and import additional chains from `wagmi/chains` or define your own custom chains. Then add them to the chains array in the wagmi configuration.

For example, to add the Polygon network:

```typescript
import { mainnet, sepolia, polygon } from 'wagmi/chains';

const config = createConfig({
  chains: [mainnet, sepolia, customChain, polygon],
  transports: {
    [mainnet.id]: http(),
    [sepolia.id]: http(),
    [customChain.id]: http(customChain.rpcUrls.default.http[0]),
    [polygon.id]: http(),
  },
  // ...
});
```

### Styling

The connect button is styled with Tailwind CSS. You can modify the styles in `app/components/ConnectButton.tsx`.

### Theme

RainbowKit theme can be customized in `app/providers.tsx`. See the [RainbowKit documentation](https://www.rainbowkit.com/docs/theming) for more options.

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [RainbowKit Documentation](https://www.rainbowkit.com/docs/introduction)
- [wagmi Documentation](https://wagmi.sh/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new) from the creators of Next.js.

Check out the [Next.js deployment documentation](https://nextjs.org/docs/deployment) for more details.
