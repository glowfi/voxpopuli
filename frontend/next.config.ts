import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
    env: {
        API_BASE_URL: process.env.API_BASE_URL,
        APP_NAME: process.env.APP_NAME,
        APP_MOTTO: process.env.APP_MOTTO
    },
    devIndicators: false
};

export default nextConfig;
