/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,

  // 👇 tells pnpm/Next to compile the shared package
  transpilePackages: ["@repo/ui"],

  images: {
    unoptimized: true,
  },

  // ⬇️  **disables TypeScript errors during `next build`**
  typescript: {
    ignoreBuildErrors: true,
  },
};

module.exports = nextConfig;
