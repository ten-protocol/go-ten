/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "export",
  // distDir should be "../api/static" in production but .next in development
  distDir: process.env.NODE_ENV === "development" ? ".next" : "../api/static",
  images: {
    unoptimized: true,
  },
};

module.exports = nextConfig;
