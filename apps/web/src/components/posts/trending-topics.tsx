'use client';
import { Separator } from '@/components/ui/separator';
import { trpc } from '@voxpopuli/trpc-client';
import { TrendingUpIcon } from 'lucide-react';
import { useTheme } from 'next-themes';
import Link from 'next/link';
import React, { useEffect, useState } from 'react';
import { MagicCard } from '../magicui/magic-card';
import { LoadingButton } from '../ui/LoadingButton';
import { topic } from './types';

const Trendingtopics = () => {
    const [alltopics, setAlltopics] = useState<topic[] | []>([]);
    const [isloading, setIsloading] = useState(false);
    const [page, setPage] = useState<number>(0);
    const [limit, setLimit] = useState<number>(6);
    const { theme } = useTheme();

    useEffect(() => {
        const fn = async () => {
            let res = await trpc.topics.getAllTopics.query({
                limit,
                skip: page * limit
            });
            // let res = await trpcClient.topics.getAllTopics.query({
            //     limit,
            //     skip: page * limit
            // });

            if (res) {
                //@ts-ignore
                setAlltopics((prev) => [...prev, ...res]);
            }
        };
        fn();
    }, [page]);

    return (
        <MagicCard
            className="flex flex-col justify-center items-center bg-muted h-fit p-6 min-w-[300px]"
            gradientColor={theme === 'dark' ? '#262626' : '#D9D9D955'}
        >
            <div className="flex flex-col items-center justify-center">
                <div className="mb-3 pointer-events-none z-10 text-wrap h-full whitespace-pre-wrap bg-gradient-to-br from-[#ff2975] from-35% to-[#00FFF1] bg-clip-text text-center text-3xl font-bold leading-none tracking-tighter text-transparent dark:drop-shadow-[0_5px_5px_rgba(0,0,0,0.8)]">
                    Trending Topics
                </div>
                <LoadingButton
                    className="mt-3"
                    loading={isloading}
                    onClick={async () => {
                        setPage((prev) => prev + 1);
                    }}
                >
                    View More
                </LoadingButton>
            </div>
            <div>
                <div className="space-y-2 mt-3">
                    {alltopics.map(({ title }, idx) => {
                        return (
                            <Link
                                key={idx}
                                href="/"
                                className="flex items-center gap-2 text-sm hover:underline"
                            >
                                <TrendingUpIcon className="h-4 w-4 text-muted-foreground" />
                                <span>{title}</span>
                            </Link>
                        );
                    })}
                </div>
            </div>
            <Separator className="my-4" />
        </MagicCard>
    );
};

export default React.memo(Trendingtopics);
