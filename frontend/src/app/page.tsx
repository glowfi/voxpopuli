'use client';
import { postApi } from '@/api/post/post';
import InfiniteScroll from '@/components/infinite-scroll';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import Posts from '@/post/Posts';
import { Loader2 } from 'lucide-react';
import React from 'react';

const Page = () => {
    const { data, next, hasMore, loading } = useInfiniteScroll<Post>(
        postApi.posts,
        10
    );
    return (
        <div className="max-h-[300px] w-full  overflow-y-auto px-10">
            <div className="flex w-full flex-col items-center  gap-3">
                <Posts posts={data} />
                <InfiniteScroll
                    hasMore={hasMore}
                    isLoading={loading}
                    next={next}
                    threshold={1}
                >
                    {hasMore && (
                        <Loader2 className="my-4 h-8 w-8 animate-spin" />
                    )}
                </InfiniteScroll>
            </div>
        </div>
    );
};

export default Page;
