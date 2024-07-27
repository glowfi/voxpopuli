'use client';
import { trpc } from '@voxpopuli/trpc-client/src/';
import React, { useState } from 'react';
import LoadingSpinner from '../loadingspinners/loadingspinner';
import InfiniteScroll from '../ui/InfiniteScroll';
import { TOTAL_POSTS_TO_FETCH } from './constants';
import CardPost from './post-cards/card-post';
import MediaViewer from './post-cards/media-viewer';
import { post } from './types';

const Posts = () => {
    const [page, setPage] = React.useState(0);
    const [loading, setLoading] = React.useState(false);
    const [hasMore, setHasMore] = React.useState(true);
    const [posts, setPosts] = useState<post[]>([]);
    const [isopen, setIsopen] = useState(false);

    const next = async () => {
        setLoading(true);

        setTimeout(async () => {
            // await new Promise((res) => setTimeout(res, 300000000));
            const data = await trpc.posts.getallposts.query({
                limit: TOTAL_POSTS_TO_FETCH,
                skip: page * TOTAL_POSTS_TO_FETCH
            });

            //@ts-ignore
            setPosts((prev) => [...prev, ...data]);
            setPage((prev) => prev + 1);

            // Usually your response will tell you if there is no more data.
            if (data.length < TOTAL_POSTS_TO_FETCH) {
                setHasMore(false);
            }
            setLoading(false);
        }, 100);
    };

    return (
        <div className="flex-1 space-y-6 p-6 flex-wrap">
            {posts ? (
                <>
                    {posts?.map((p, idx) => {
                        return (
                            <div
                                className="relative w-full overflow-hidden rounded-lg md:shadow-xl m-0 p-0"
                                key={idx}
                            >
                                <CardPost
                                    p={p}
                                    open={isopen}
                                    setIsopen={setIsopen}
                                />
                                <MediaViewer
                                    open={isopen}
                                    setIsopen={setIsopen}
                                />
                            </div>
                        );
                    })}
                </>
            ) : (
                ''
            )}
            <InfiniteScroll
                hasMore={hasMore}
                isLoading={loading}
                next={next}
                threshold={1}
            >
                {hasMore && (
                    <div className="m-auto flex justify-center items-center justify-items-center">
                        <LoadingSpinner name="more posts" />
                    </div>
                )}
            </InfiniteScroll>
        </div>
    );
};

export default React.memo(Posts);
