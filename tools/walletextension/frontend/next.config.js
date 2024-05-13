/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "export",
  // distDir should be "../api/static" in production but .next in development
  distDir: process.env.NODE_ENV === "development" ? ".next" : "../api/static",
  images: {
    unoptimized: true,
  },
  // base path for static files should be "" in development but "/static" in production
  basePath: process.env.NODE_ENV === "development" ? "" : "/static",
};

module.exports = nextConfig;
