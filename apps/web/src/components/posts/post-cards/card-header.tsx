import { Badge } from '@/components/ui/badge';
import Image from 'next/image';
import Link from 'next/link';
import { Avatar, AvatarFallback, AvatarImage } from '../../ui/avatar';
import { post } from '../types';

const PostCardHeader = ({ post }: { post: post }) => {
    return (
        <div className="flex flex-row justify-start items-center px-4 py-3 bg-muted gap-3">
            <div className="flex justify-center items-center gap-3">
                <div className="flex gap-1 justify-center items-center">
                    <span className="text-yellow-600 font-bold">By</span>
                    <Avatar className="w-10 h-10 border-2 hover:transition hover:opacity-75 hover:cursor-pointer bg-white">
                        <AvatarImage src={post?.author.avatarImg} />
                        <AvatarFallback>
                            {post?.author.username.charAt(0)}
                        </AvatarFallback>
                    </Avatar>

                    <Link
                        className="text-sm font-semibold hover:underline dark:text-white text-wrap"
                        href={`/u/${post?.author?.username}`}
                    >
                        u/{post?.author?.username}
                    </Link>
                </div>

                <div className="hidden md:flex gap-1 justify-center items-center">
                    <span className="text-red-500 font-bold mr-1">in</span>
                    <Avatar className="w-12 h-12 border-2 border-background hover:transition hover:opacity-75 hover:cursor-pointer">
                        {/* @ts-ignore */}
                        <AvatarImage src={post?.voxsphere?.[0]?.logoURL} />
                        <AvatarFallback>v/</AvatarFallback>
                    </Avatar>
                    <Link
                        className="text-sm font-semibold hover:underline"
                        href={`v/${post?.voxsphere?.[0]?.name}`}
                    >
                        v/{post?.voxsphere?.[0]?.name}
                    </Link>
                </div>

                <div className="hidden md:flex gap-1 justify-center items-center">
                    {post?.postflair?.[0]?.title && (
                        <Badge
                            style={{
                                backgroundColor:
                                    post?.postflair?.[0].colorHex || ''
                            }}
                        >
                            {post?.postflair?.[0]?.title}
                        </Badge>
                    )}
                </div>
                <div className="hidden md:flex gap-1 justify-center items-center">
                    {post?.awards?.map(({ imageLink }, idx) => {
                        return (
                            <Image
                                key={idx}
                                alt="Not Found"
                                src={imageLink}
                                height={20}
                                width={20}
                            />
                        );
                    })}
                </div>
                <div className="hidden sm:flex gap-1 justify-center items-center">
                    <span className="font-bold">·</span>
                    <span className="font-bold"> {post?.createdatHuman}</span>
                </div>
            </div>
        </div>
    );
};

export default PostCardHeader;
