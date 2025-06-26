import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",         // ONLY proxy /api routes
        destination: "http://localhost:8080/:path*", // Your Go backend
      },
    ];
  },
};

export default nextConfig;
