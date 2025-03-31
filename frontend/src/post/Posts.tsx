'use client';
import React from 'react';
import PostCard from './PostCard';
import { Post } from './post';

interface PostsProps {
    posts: Post[];
}

const Posts: React.FunctionComponent<PostsProps> = ({ posts }) => {
    return (
        <div className="flex flex-col gap-3">
            {posts.map((p, idx) => {
                return <PostCard post={p} key={idx} />;
            })}
        </div>
    );
};

export default Posts;
