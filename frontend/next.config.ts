import type { NextConfig } from "next";
// import { NextRequest, NextResponse } from "next/server";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/:path*",
        destination: "http://localhost:8080/:path*", // Proxy to Go backend
      },
    ];
  },
};

export default nextConfig;