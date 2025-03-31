'use client';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { formatDistanceToNow } from 'date-fns';
import {
    Award,
    MessageSquare,
    Share2,
    ThumbsDown,
    ThumbsUp
} from 'lucide-react';
import { FunctionComponent, useState } from 'react';
import {
    GalleryMetadata,
    GifMetadata,
    ImageMetadata,
    LinkType,
    MediaType,
    Post,
    Video
} from './post';
import { Separator } from '@/components/ui/separator';
import Gallery from './content/gallery';
import ImageViz from './content/image';
import Gif from './content/gif';
import Text from './content/text';
import VideoViz from './content/video';
import LinkViz from './content/link';
import Link from 'next/link';

interface PostCardProps {
    post: Post;
}

const PostCard: FunctionComponent<PostCardProps> = ({ post }) => {
    const [votes, setVotes] = useState(post.ups);
    const [userVote, setUserVote] = useState<'up' | 'down' | null>(null);

    const handleVote = (direction: 'up' | 'down') => {
        if (userVote === direction) {
            // Undo vote
            setVotes(direction === 'up' ? votes - 1 : votes + 1);
            setUserVote(null);
        } else {
            // Change vote
            if (userVote === 'up' && direction === 'down') {
                setVotes(votes - 2);
            } else if (userVote === 'down' && direction === 'up') {
                setVotes(votes + 2);
            } else {
                // New vote
                setVotes(direction === 'up' ? votes + 1 : votes - 1);
            }
            setUserVote(direction);
        }
    };

    const renderContent = () => {
        switch (post.media_type) {
            case MediaType.Image:
                return <ImageViz data={post.medias as ImageMetadata[]} />;
            case MediaType.Gif:
                return <Gif data={post.medias as GifMetadata[]} />;
            case MediaType.Gallery:
                return <Gallery data={post.medias as GalleryMetadata[]} />;
            case MediaType.Video:
                return <VideoViz data={post?.medias?.[0] as Video} />;
            case MediaType.Link:
                return (
                    <LinkViz
                        data={post?.medias?.[0] as LinkType}
                        isFullPage={false}
                    />
                );
            case MediaType.Text:
                return <Text data={{ text: post.text }} />;
            default:
                return null;
        }
    };

    return (
        <div className="bg-card rounded-lg shadow-md overflow-hidden transition-all hover:shadow-lg border">
            {/* Post Header */}
            <div className="p-4 flex items-start gap-3">
                <Avatar className="h-10 w-10 border">
                    <AvatarImage
                        src={`/placeholder.svg?height=40&width=40&text=${post.author}`}
                    />
                    <AvatarFallback>{post.author}</AvatarFallback>
                </Avatar>

                <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                        <Badge
                            variant="outline"
                            className="bg-primary/10 hover:bg-primary/20 font-bold hover:cursor-pointer"
                        >
                            {post.voxsphere}
                        </Badge>
                        <span className="text-sm text-slate-500 dark:text-slate-400">
                            Posted by{' '}
                            <Link href="#" className="hover:underline">
                                {' '}
                                u/{post.author}
                            </Link>{' '}
                            •{' '}
                            {formatDistanceToNow(post.created_at, {
                                addSuffix: true
                            })}
                        </span>
                    </div>

                    <Link href="#" className="hover:underline">
                        <h2 className="text-lg font-bold mt-1 text-slate-900 dark:text-slate-100">
                            {post.title}
                        </h2>
                    </Link>
                </div>
            </div>

            {/* Post Content */}
            <div className="px-4 pb-2">{renderContent()}</div>

            <Separator />

            {/* Post Footer */}
            <div className="p-2 flex items-center justify-between flex-wrap gap-2">
                <div className="flex items-center gap-1 sm:gap-2">
                    <Button
                        variant="ghost"
                        size="sm"
                        className={cn(
                            'text-slate-600 dark:text-slate-300',
                            userVote === 'up' &&
                                'text-green-500 dark:text-green-400'
                        )}
                        onClick={() => handleVote('up')}
                    >
                        <ThumbsUp className="h-4 w-4 mr-1" />
                        <span className="text-xs sm:text-sm">
                            {votes > 0 ? `+${votes}` : votes}
                        </span>
                    </Button>

                    <Button
                        variant="ghost"
                        size="sm"
                        className={cn(
                            'text-slate-600 dark:text-slate-300',
                            userVote === 'down' &&
                                'text-red-500 dark:text-red-400'
                        )}
                        onClick={() => handleVote('down')}
                    >
                        <ThumbsDown className="h-4 w-4" />
                    </Button>
                </div>

                <div className="flex items-center gap-1 sm:gap-3">
                    {post.num_comments > 0 && (
                        <Button
                            variant="ghost"
                            size="sm"
                            className="text-slate-600 dark:text-slate-300"
                        >
                            <MessageSquare className="h-4 w-4 mr-1" />
                            <span className="text-xs sm:text-sm">
                                {post.num_comments} comments
                            </span>
                        </Button>
                    )}

                    <Button
                        variant="ghost"
                        size="sm"
                        className="text-slate-600 dark:text-slate-300"
                    >
                        <Share2 className="h-4 w-4 mr-1" />
                        <span className="text-xs sm:text-sm">Share</span>
                    </Button>

                    {post.num_awards > 0 && (
                        <Button
                            variant="ghost"
                            size="sm"
                            className="text-slate-600 dark:text-slate-300"
                        >
                            <Award className="h-4 w-4 mr-1" />
                            <span className="text-xs sm:text-sm">
                                {post.num_awards}
                            </span>
                        </Button>
                    )}
                </div>
            </div>
        </div>
    );
};

export default PostCard;
