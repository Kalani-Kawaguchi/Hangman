import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",         // ONLY proxy /api routes
        destination: `${process.env.REACT_APP_API_URL}/:path*`, // Your Go backend
      },
    ];
  },
};

export default nextConfig;
