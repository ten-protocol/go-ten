import { socialLinks } from "@repo/ui/lib/constants";

export const siteMetadata = {
  companyName: "Tenscan",
  metaTitle: "Tenscan - Real-Time Blockchain Explorer",
  description:
    "Tenscan allows you to explore and search the TEN blockchain. View transaction history and smart contract interactions without compromising sensitive data.",
  keywords:
    "blockchain explorer, real-time blockchain data, dapps, l2, encryption, layer2, crypto transactions, blockchain analysis, Tenscan, block explorer,TEN, TEN Protocol, TEN chain, TEN network, TEN blockchain, TEN explorer",
  siteUrl: "https://tenscan.io",
  siteLogo: `/assets/images/cover.png`,
  siteLogoSquare: `/assets/images/cover.png`,
  email: socialLinks.email,
  twitter: socialLinks.twitter,
  twitterHandle: socialLinks.twitterHandle,
  github: socialLinks.github,

  rollups: {
    title: "Tenscan | Rollups",
    description: "View the latest rollups on the TEN network.",
    canonicalUrl: "https://tenscan.io/rollups",
    ogImageUrl: "/assets/images/cover.png",
    ogTwitterImage: "/assets/images/cover.png",
    ogType: "website",
  },

  blocks: {
    title: "Tenscan | Blocks",
    description:
      "Analyze blocks on TEN in real-time. Access block height, timestamp, transaction count, and miner details instantly",
    canonicalUrl: "https://tenscan.io/blocks",
    ogImageUrl: "/assets/images/cover.png",
    ogTwitterImage: "/assets/images/cover.png",
    ogType: "website",
  },

  batches: {
    title: "Tenscan | Batches",
    description:
      "Explore the history of batches on TEN. Access detailed batch lists, transaction summaries, and batch metadata",
    canonicalUrl: "https://tenscan.io/batches",
    ogImageUrl: "/assets/images/cover.png",
    ogTwitterImage: "/assets/images/cover.png",
    ogType: "website",
  },

  transactions: {
    title: "Tenscan | Public Transactions",
    description:
      "Explore all public transactions on TEN. View essential details including batch, batch age, transaction hash, and finality status",
    canonicalUrl: "https://tenscan.io/transactions",
    ogImageUrl: "/assets/images/cover.png",
    ogTwitterImage: "/assets/images/cover.png",
    ogType: "website",
  },

  personal: {
    title: "Tenscan | Personal Transactions",
    description:
      "Access detailed information about your transactions on TEN. View sender, receiver, amount, confirmation status and more",
    canonicalUrl: "https://tenscan.io/personal",
    ogImageUrl: "/assets/images/cover.png",
    ogTwitterImage: "/assets/images/cover.png",
    ogType: "website",
  },
};
