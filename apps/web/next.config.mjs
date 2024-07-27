/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: false,
    env: {
        SERVER_URL: process.env.SERVER_URL,
        APP_NAME: process.env.APP_NAME
    },
    images: {
        remotePatterns: [
            {
                hostname: 'res.cloudinary.com'
            },
            {
                hostname: 'ik.imagekit.io'
            },
            {
                hostname: 'external-preview.redd.it'
            },
            {
                hostname: 'preview.redd.it'
            },
            {
                hostname: 'styles.redditmedia.com'
            },
            {
                hostname: 'www.redditstatic.com'
            },
            { hostname: 'i.redd.it' },
            { hostname: 'robohash.org' },
            { hostname: 'avatars.githubusercontent.com' }
        ]
    }
};

export default nextConfig;
