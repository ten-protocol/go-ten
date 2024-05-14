/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  images: {
    unoptimized: true,
  },
  async redirects() {
    const destinationUrl = process.env.NEXT_PUBLIC_API_GATEWAY_URL || 'https://testnet.ten.xyz';
    return [
      {
        source: '/v1/:path*',
        destination: `${destinationUrl}/v1/:path*`,
        permanent: true,
        has: [
          {
            type: 'header',
            key: 'x-forwarded-proto',
            value: '(https|http)',
          },
        ],
      },
    ]
  },
};

module.exports = nextConfig;
