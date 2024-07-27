'use client';
import { trpcClient as trpc } from './client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { httpBatchLink } from '@trpc/client';
import React, { useState } from 'react';
import 'dotenv/config';
import { configDotenv } from 'dotenv';

configDotenv({
    path: '../.env'
});

export const Provider = ({ children }: { children: React.ReactNode }) => {
    console.log('CHECK', process.env.SERVER_URL);
    const [queryClient] = useState(() => new QueryClient());
    const [trpcClient] = useState(() =>
        trpc.createClient({
            links: [
                httpBatchLink({
                    url: process.env.SERVER_URL as string
                    // url: 'http://localhost:8080/trpc'
                })
            ]
        })
    );

    return (
        <trpc.Provider client={trpcClient} queryClient={queryClient}>
            <QueryClientProvider client={queryClient}>
                {children}
            </QueryClientProvider>
        </trpc.Provider>
    );
};
