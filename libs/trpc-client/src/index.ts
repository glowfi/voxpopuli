import type { AppRouter } from '@voxpopuli/trpc-server/routers';
import { createTRPCProxyClient, httpBatchLink } from '@trpc/client';
import 'dotenv/config';

export const trpc = createTRPCProxyClient<AppRouter>({
    links: [
        httpBatchLink({
            // url: process.env.SERVER_URL as string
            url: 'http://localhost:8080/trpc'
        })
    ]
});
