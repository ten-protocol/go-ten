/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "export",
  // distDir is where the static files are generated
  // on dev mode, it is generated in .next folder in the root directory
  // comment this during dev mode
  distDir: "../api/static",
};

module.exports = nextConfig;
