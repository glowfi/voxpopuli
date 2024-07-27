'use client';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useTheme } from 'next-themes';
import Link from 'next/link';
import React, { useEffect, useState } from 'react';
import LoadingSpinner from '../loadingspinners/loadingspinner';
import { MagicCard } from '../magicui/magic-card';
import { LoadingButton } from '../ui/LoadingButton';
import { voxsphere } from './types';
import { trpc } from '@voxpopuli/trpc-client/src/';

const PopularCommunities = () => {
    const [PopularCommunities, setPopularCommunities] = useState<
        null | [] | voxsphere[]
    >([]);
    const [isloading, setIsloading] = useState(true);
    const [limit, setLimit] = useState<number>(6);
    const { theme } = useTheme();

    useEffect(() => {
        const fn = async () => {
            try {
                // await new Promise((res) => setTimeout(res, 3000));
                // let data = await trpcClient.voxspheres.getTopVoxspheres.query({
                //     limit
                // });
                let data = await trpc.voxspheres.getTopVoxspheres.query({
                    limit
                });
                console.log(data, 'GOT BACK');
                setPopularCommunities(data);
                setIsloading(false);
            } catch {
                setIsloading(false);
            }
            setIsloading(false);
        };
        fn();
    }, [limit]);

    return (
        <MagicCard
            className="flex flex-col justify-center items-center bg-muted h-fit p-6 min-w-[300px]"
            gradientColor={theme === 'dark' ? '#262626' : '#D9D9D955'}
        >
            <div className="flex flex-col items-center justify-between">
                <div className="mb-3 pointer-events-none z-10 text-wrap h-full whitespace-pre-wrap bg-gradient-to-br from-[#ff2975] from-35% to-[#00FFF1] bg-clip-text text-center text-3xl font-bold leading-none tracking-tighter text-transparent dark:drop-shadow-[0_5px_5px_rgba(0,0,0,0.8)]">
                    Popular Communities
                </div>
                <LoadingButton
                    loading={isloading}
                    onClick={() => {
                        setIsloading(true);
                        setLimit((prev) => prev + 2);
                    }}
                    className="mb-3"
                >
                    View More
                </LoadingButton>
            </div>
            <div>
                <div className="space-y-2">
                    {PopularCommunities?.map((p, idx) => {
                        return (
                            <Link
                                key={idx}
                                href="#"
                                className="flex items-center gap-2 text-sm hover:underline"
                            >
                                <Avatar className="h-10 w-10">
                                    {/*@ts-ignore */}
                                    <AvatarImage src={p?.logoURL} />
                                    <AvatarFallback>
                                        {p?.name.charAt(0)}
                                    </AvatarFallback>
                                </Avatar>
                                <div className="flex flex-wrap gap-3">
                                    <span>v/{p?.name}</span>
                                    <span className="font-bold">
                                        {p?.totalmembersHuman}
                                    </span>
                                </div>
                            </Link>
                        );
                    })}
                </div>
                {isloading ? (
                    <div className="flex justify-center items-center mx-auto">
                        <LoadingSpinner name="top communities" />
                    </div>
                ) : (
                    ''
                )}
            </div>
        </MagicCard>
    );
};

export default React.memo(PopularCommunities);
