import AnimatedGradientText from '@/components/magicui/animated-gradient-text';
import AvatarCircles from '@/components/magicui/avatar-circles';
import { cn } from '@/lib/utils';
import {
    AwardIcon,
    ChevronDownIcon,
    ChevronUpIcon,
    Medal,
    Share
} from 'lucide-react';
import { post } from '../types';
import { getAvatarsArray } from '../utils';

const PostCardFooter = ({ p }: { p: post }) => {
    return (
        <div className="w-full flex items-center justify-start md:justify-between px-4 py-3 bg-muted gap-3">
            <div className="flex justify-center items-center gap-2 text-sm text-muted-foreground">
                <div className="hidden md:flex flex-wrap gap-1 hover:transition hover:opacity-75 hover:cursor-pointer justify-center items-center">
                    <Medal className="w-6 h-6" />
                    <span>Give Award</span>
                </div>
                <div className="hidden md:flex gap-1 hover:transition hover:opacity-75 hover:cursor-pointer justify-center items-center">
                    <Share className="w-6 h-6" />
                    <span>Share</span>
                </div>
            </div>
            <div className="flex justify-around md:justify-center items-center text-sm text-muted-foreground mx-auto sm:mx-0 gap-3">
                {p?.numComments > 0 && (
                    <div className="hidden sm:flex justify-center items-center gap-1 hover:transition hover:opacity-75 hover:cursor-pointer border">
                        <AvatarCircles
                            numPeople={p?.numComments}
                            avatarUrls={getAvatarsArray()}
                        />
                    </div>
                )}
                {p?.awards?.length > 0 && (
                    <div className="flex gap-1 hover:transition hover:opacity-75 hover:cursor-pointer justify-center items-center">
                        <AwardIcon className="w-6 h-6" />
                        <span>{p?.awards?.length}</span>
                    </div>
                )}

                <div className="flex gap-3 justify-center items-center">
                    <AnimatedGradientText className="hover:cursor-pointer">
                        <span
                            className={cn(
                                `inline animate-gradient bg-gradient-to-r from-[#ffaa40] via-[#9c40ff] to-[#ffaa40] bg-[length:var(--bg-size)_100%] bg-clip-text text-transparent`
                            )}
                        >
                            Uv
                        </span>
                        <ChevronUpIcon className="ml-1 size-3 transition-transform duration-300 ease-in-out group-hover:translate-x-1" />
                    </AnimatedGradientText>
                    <span>{p?.ups}</span>
                    <AnimatedGradientText className="hover:cursor-pointer">
                        <span
                            className={cn(
                                `inline animate-gradient bg-gradient-to-r from-[#ffaa40] via-[#9c40ff] to-[#ffaa40] bg-[length:var(--bg-size)_100%] bg-clip-text text-transparent`
                            )}
                        >
                            Dv
                        </span>
                        {/* <Loader2 className="animate-spin w-3 h-3 ml-3" /> */}
                        <ChevronDownIcon className="ml-1 size-3 transition-transform duration-300 ease-in-out group-hover:translate-x-1" />
                    </AnimatedGradientText>
                </div>
            </div>
        </div>
    );
};

export default PostCardFooter;
