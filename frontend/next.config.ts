import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
    /* config options here */
    env: {
        API_ENDPOINT: process.env.API_ENDPOINT
    }
};

export default nextConfig;
