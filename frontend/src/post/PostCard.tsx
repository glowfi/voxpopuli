import React, { FunctionComponent } from 'react';

interface PostCardProps {
    post: Post;
}

const PostCard: FunctionComponent<PostCardProps> = ({ post }) => {
    return <div>{post.title}</div>;
};

export default PostCard;
