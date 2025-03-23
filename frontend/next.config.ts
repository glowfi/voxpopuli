import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
    env: {
        API_ENDPOINT: process.env.API_ENDPOINT,
        APP_NAME: process.env.APP_NAME,
        APP_MOTTO: process.env.APP_MOTTO
    },
    devIndicators: false
};

export default nextConfig;
