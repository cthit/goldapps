/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  rewrites: async () => {
    return [
      {
        source: "/api/:path*",
        destination: `${String(process.env.GOLDAPPS_URL)}/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
