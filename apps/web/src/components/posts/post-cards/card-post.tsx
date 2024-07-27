'use client';
import { BorderBeam } from '@/components/magicui/border-beam';
import { MagicCard } from '@/components/magicui/magic-card';
import DOMPurify from 'dompurify';
import parse from 'html-react-parser';
import { useTheme } from 'next-themes';
import Link from 'next/link';
import React from 'react';
import HandleMedia from '../handle-media';
import { post } from '../types';
import PostCardFooter from './card-footer';
import PostCardHeader from './card-header';

interface props {
    p: post;
    open: boolean;
    setIsopen: React.Dispatch<React.SetStateAction<boolean>>;
}

const CardPost = ({ p }: props) => {
    const { theme } = useTheme();
    return (
        <MagicCard
            className="flex flex-col justify-center gap-3"
            gradientColor={theme === 'dark' ? '#262626' : '#D9D9D955'}
        >
            <PostCardHeader post={p} />
            <div>
                <div className="flex flex-col justify-center items-center gap-4">
                    <div className="flex flex-col flex-wrap mt-3">
                        <Link
                            className="scroll-m-20 text-xl font-semibold tracking-tight first:mt-0 hover:underline hover:cursor-pointer hover:opacity-80 text-center"
                            href={`/thepost/${p.id}`}
                        >
                            {p?.title}
                        </Link>
                    </div>
                    <div className="flex flex-col flex-wrap justify-center items-center gap-3">
                        {p?.mediaContent?.[0] && (
                            <HandleMedia
                                {...{
                                    mediadata: p?.mediaContent?.[0],
                                    over18: p?.over18,
                                    spoiler: p?.spoiler
                                }}
                            />
                        )}
                    </div>
                    <p className="line-clamp-6 text-wrap max-w-[150px] md:max-w-[450px] text-muted-foreground animate-pulse mb-3 text-center">
                        {p?.textHTML && parse(DOMPurify.sanitize(p?.textHTML))}
                    </p>
                </div>
            </div>
            <PostCardFooter p={p} />
            <BorderBeam size={250} duration={12} delay={9} />
        </MagicCard>
    );
};

export default CardPost;
