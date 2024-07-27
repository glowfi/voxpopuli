import type { AppRouter } from '@voxpopuli/trpc-server/routers';
import { createTRPCProxyClient, httpBatchLink } from '@trpc/client';
import 'dotenv/config';
import { configDotenv } from 'dotenv';

configDotenv({
    path: '../.env'
});

console.log('CHECK', process.env.SERVER_URL);
export const trpc = createTRPCProxyClient<AppRouter>({
    links: [
        httpBatchLink({
            // url: process.env.SERVER_URL as string
            url: 'https://voxpopuli.onrender.com/trpc'
        })
    ]
});
