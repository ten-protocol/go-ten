/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "export",
  // distDir should be "../api/static" in production but .next in development
  distDir: process.env.NODE_ENV === "production" ? "../api/static" : ".next",
};

module.exports = nextConfig;
