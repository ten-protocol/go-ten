/** @type {import('next').NextConfig} */
const path = require("path");

const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ["@repo/ui"],
  // TODO: temp fix until a proper solution to the build errors is implemented.
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  webpack: (config) => {
    const reactPkgDir = path.dirname(require.resolve("react/package.json"));
    const reactDomPkgDir = path.dirname(require.resolve("react-dom/package.json"));

    config.resolve = config.resolve || {};
    config.resolve.alias = {
      ...(config.resolve.alias || {}),
      react: reactPkgDir,
      "react-dom": reactDomPkgDir,
    };
    return config;
  },
};

module.exports = nextConfig;
