import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
    /* config options here */
    env: {
        API_ENDPOINT: process.env.API_ENDPOINT,
        APP_NAME: process.env.APP_NAME,
        APP_MOTTO: process.env.APP_MOTTO
    }
};

export default nextConfig;
