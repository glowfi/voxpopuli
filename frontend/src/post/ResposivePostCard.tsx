'use client';

import { useState } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
    Card,
    CardContent,
    CardFooter,
    CardHeader
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
    ArrowDown,
    ArrowUp,
    Award,
    BookmarkPlus,
    MessageSquare,
    MoreHorizontal,
    Share2
} from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';

export function ResponsivePostCard() {
    const [votes, setVotes] = useState(1452);
    const [userVote, setUserVote] = useState<'up' | 'down' | null>(null);

    const handleVote = (vote: 'up' | 'down') => {
        if (userVote === vote) {
            setVotes(votes + (vote === 'up' ? -1 : 1));
            setUserVote(null);
        } else {
            setVotes(
                votes +
                    (vote === 'up'
                        ? userVote === 'down'
                            ? 2
                            : 1
                        : userVote === 'up'
                          ? -2
                          : -1)
            );
            setUserVote(vote);
        }
    };

    return (
        <Card className="w-full max-w-full mx-auto overflow-hidden border-t-0 rounded-t-none sm:border-t sm:rounded-t mb-3">
            {/* Desktop and Tablet View: Side vote buttons */}
            <div className="hidden sm:flex sm:flex-row">
                <div className="flex-1">
                    <CardHeader className="p-3 pb-0 space-y-0">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center flex-wrap gap-1 text-sm">
                                <Avatar className="w-5 h-5">
                                    <AvatarImage
                                        src="/placeholder.svg?height=20&width=20"
                                        alt="r/photography"
                                    />
                                    <AvatarFallback>P</AvatarFallback>
                                </Avatar>
                                <Link
                                    href="#"
                                    className="font-medium hover:underline"
                                >
                                    r/photography
                                </Link>
                                <span className="text-muted-foreground">
                                    • Posted by
                                </span>
                                <Link
                                    href="#"
                                    className="text-muted-foreground hover:underline"
                                >
                                    u/photomaster
                                </Link>
                                <span className="text-muted-foreground">
                                    5h
                                </span>
                            </div>
                            <div className="flex items-center">
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    className="h-8 w-8"
                                >
                                    <BookmarkPlus className="w-4 h-4" />
                                </Button>
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            className="h-8 w-8"
                                        >
                                            <MoreHorizontal className="w-4 h-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end">
                                        <DropdownMenuItem>
                                            Save
                                        </DropdownMenuItem>
                                        <DropdownMenuItem>
                                            Hide
                                        </DropdownMenuItem>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem>
                                            Report
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </div>
                        </div>
                        <h3 className="pt-2 text-base font-medium md:text-lg">
                            Sunset over the mountains - captured this amazing
                            view last night
                        </h3>
                    </CardHeader>
                    <CardContent className="p-3">
                        <div className="relative overflow-hidden rounded-md aspect-video">
                            <Image
                                src="/placeholder.svg?height=400&width=600"
                                alt="Sunset over mountains"
                                fill
                                className="object-cover"
                            />
                        </div>
                    </CardContent>
                    <CardFooter className="flex items-center justify-between p-3 pt-0">
                        <div className="flex items-center space-x-4">
                            <div className="flex items-center">
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    className={`h-8 w-8 ${userVote === 'up' ? 'text-orange-500' : ''}`}
                                    onClick={() => handleVote('up')}
                                >
                                    <ArrowUp className="w-4 h-4" />
                                </Button>
                                <span className="text-sm font-medium">
                                    {votes}
                                </span>
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    className={`h-8 w-8 ${userVote === 'down' ? 'text-blue-500' : ''}`}
                                    onClick={() => handleVote('down')}
                                >
                                    <ArrowDown className="w-4 h-4" />
                                </Button>
                            </div>
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 space-x-1"
                            >
                                <MessageSquare className="w-4 h-4" />
                                <span className="text-xs md:text-sm">
                                    243 comments
                                </span>
                            </Button>
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 space-x-1"
                            >
                                <Share2 className="w-4 h-4" />
                                <span className="text-xs md:text-sm">
                                    Share
                                </span>
                            </Button>
                        </div>
                        <div className="flex items-center">
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 space-x-1"
                            >
                                <Award className="w-4 h-4" />
                                <span className="text-xs md:text-sm">
                                    Award
                                </span>
                            </Button>

                            {/* <Badge */}
                            {/*     variant="outline" */}
                            {/*     className="flex items-center space-x-1 border-amber-300" */}
                            {/* > */}
                            {/*     <Award className="w-3 h-3 text-amber-400 fill-amber-400" /> */}
                            {/*     <span className="text-xs">3</span> */}
                            {/* </Badge> */}
                        </div>
                    </CardFooter>
                </div>
            </div>

            {/* Mobile View: Top-to-bottom layout */}
            <div className="sm:hidden">
                <CardHeader className="p-3 pb-0 space-y-0">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center flex-wrap gap-1 text-xs">
                            <Avatar className="w-4 h-4">
                                <AvatarImage
                                    src="/placeholder.svg?height=20&width=20"
                                    alt="r/photography"
                                />
                                <AvatarFallback>P</AvatarFallback>
                            </Avatar>
                            <Link
                                href="#"
                                className="font-medium hover:underline"
                            >
                                r/photography
                            </Link>
                            <span className="text-muted-foreground">•</span>
                            <Link
                                href="#"
                                className="text-muted-foreground hover:underline"
                            >
                                u/photomaster
                            </Link>
                            <span className="text-muted-foreground">5h</span>
                        </div>
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    className="h-7 w-7"
                                >
                                    <MoreHorizontal className="w-4 h-4" />
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                                <DropdownMenuItem>Save</DropdownMenuItem>
                                <DropdownMenuItem>Hide</DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>Report</DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>
                    <h3 className="pt-1 text-sm font-medium">
                        Sunset over the mountains - captured this amazing view
                        last night
                    </h3>
                </CardHeader>
                <CardContent className="p-3">
                    <div className="relative overflow-hidden rounded-md aspect-video">
                        <Image
                            src="/placeholder.svg?height=400&width=600"
                            alt="Sunset over mountains"
                            fill
                            className="object-cover"
                        />
                    </div>
                </CardContent>
                <CardFooter className="flex flex-col gap-2 p-2">
                    <div className="flex items-center justify-between w-full">
                        <div className="flex items-center">
                            <Button
                                variant="ghost"
                                size="icon"
                                className={`h-8 w-8 ${userVote === 'up' ? 'text-orange-500' : ''}`}
                                onClick={() => handleVote('up')}
                            >
                                <ArrowUp className="w-4 h-4" />
                            </Button>
                            <span className="text-sm font-medium">{votes}</span>
                            <Button
                                variant="ghost"
                                size="icon"
                                className={`h-8 w-8 ${userVote === 'down' ? 'text-blue-500' : ''}`}
                                onClick={() => handleVote('down')}
                            >
                                <ArrowDown className="w-4 h-4" />
                            </Button>
                        </div>
                        <div className="flex items-center gap-2">
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 px-2"
                            >
                                <MessageSquare className="w-4 h-4 mr-1" />
                                <span className="text-xs">243</span>
                            </Button>
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 px-2"
                            >
                                <Share2 className="w-4 h-4" />
                            </Button>
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 px-2"
                            >
                                <Award className="w-4 h-4" />
                            </Button>
                            {/* <Badge */}
                            {/*     variant="outline" */}
                            {/*     className="flex items-center space-x-1 border-amber-300 h-6" */}
                            {/* > */}
                            {/*     <Award className="w-3 h-3 text-amber-400 fill-amber-400" /> */}
                            {/*     <span className="text-xs">3</span> */}
                            {/* </Badge> */}
                        </div>
                    </div>
                </CardFooter>
            </div>
        </Card>
    );
}
