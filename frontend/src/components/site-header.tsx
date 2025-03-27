'use client';

import { Search, Shield, SidebarIcon, X } from 'lucide-react';

import { SearchForm } from '@/components/search-form';
import { Button } from '@/components/ui/button';
import { useSidebar } from '@/components/ui/sidebar';
import Link from 'next/link';
import Image from 'next/image';
import { Input } from './ui/input';
import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useIsMobile } from '@/hooks/use-mobile';
import { Avatar } from './ui/avatar';
import { AvatarFallback, AvatarImage } from '@radix-ui/react-avatar';

export function SiteHeader() {
    const { toggleSidebar } = useSidebar();
    const [isLoggedIn, setIsLoggedIn] = useState(true);
    const [notificationCount, setNotificationCount] = useState(3);
    const [unreadMessages, setUnreadMessages] = useState(2);
    const [showNotifications, setShowNotifications] = useState(false);
    const [showChat, setShowChat] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');
    const [isSearching, setIsSearching] = useState(false);
    const [searchResults, setSearchResults] = useState<{
        subreddits: Array<{ name: string; members: number; icon: string }>;
        users: Array<{ username: string; karma: number; avatar: string }>;
    }>({
        subreddits: [],
        users: []
    });
    const notificationRef = useRef<HTMLDivElement>(null);
    const chatRef = useRef<HTMLDivElement>(null);
    const searchRef = useRef<HTMLDivElement>(null);
    const searchInputRef = useRef<HTMLInputElement>(null);
    const router = useRouter();
    const isMobile = useIsMobile();

    // Sample search data
    const allSubreddits = [
        {
            name: 'hiking',
            members: 2.4,
            icon: '/placeholder.svg?height=32&width=32'
        },
        {
            name: 'programming',
            members: 5.2,
            icon: '/placeholder.svg?height=32&width=32'
        },
        {
            name: 'askreddit',
            members: 42.0,
            icon: '/placeholder.svg?height=32&width=32'
        },
        {
            name: 'gaming',
            members: 35.7,
            icon: '/placeholder.svg?height=32&width=32'
        },
        {
            name: 'movies',
            members: 28.9,
            icon: '/placeholder.svg?height=32&width=32'
        }
    ];

    const allUsers = [
        {
            username: 'mountain_lover',
            karma: 24689,
            avatar: '/placeholder.svg?height=32&width=32'
        },
        {
            username: 'code_wizard',
            karma: 56432,
            avatar: '/placeholder.svg?height=32&width=32'
        },
        {
            username: 'film_buff',
            karma: 12543,
            avatar: '/placeholder.svg?height=32&width=32'
        },
        {
            username: 'hiking_parent',
            karma: 8765,
            avatar: '/placeholder.svg?height=32&width=32'
        },
        {
            username: 'travel_enthusiast',
            karma: 34567,
            avatar: '/placeholder.svg?height=32&width=32'
        }
    ];

    // Filter results based on search query
    useEffect(() => {
        if (searchQuery.trim() === '') {
            setSearchResults({
                subreddits: [],
                users: []
            });
            return;
        }

        const query = searchQuery.toLowerCase();
        const filteredSubreddits = allSubreddits.filter((subreddit) =>
            subreddit.name.toLowerCase().includes(query)
        );

        const filteredUsers = allUsers.filter((user) =>
            user.username.toLowerCase().includes(query)
        );

        setSearchResults({
            subreddits: filteredSubreddits,
            users: filteredUsers
        });
    }, [searchQuery]);

    // Close panels when clicking outside
    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (
                notificationRef.current &&
                !notificationRef.current.contains(event.target as Node)
            ) {
                setShowNotifications(false);
            }
            if (
                chatRef.current &&
                !chatRef.current.contains(event.target as Node)
            ) {
                setShowChat(false);
            }
            if (
                searchRef.current &&
                !searchRef.current.contains(event.target as Node) &&
                searchInputRef.current &&
                !searchInputRef.current.contains(event.target as Node)
            ) {
                setIsSearching(false);
            }
        }

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    const handleCreatePost = () => {
        router.push('/submit');
    };

    const handleNotificationClick = () => {
        setShowNotifications(!showNotifications);
        if (!showNotifications) {
            setShowChat(false); // Close chat if opening notifications
            setIsSearching(false); // Close search if opening notifications
        }
    };

    const handleChatClick = () => {
        setShowChat(!showChat);
        if (!showChat) {
            setShowNotifications(false); // Close notifications if opening chat
            setIsSearching(false); // Close search if opening chat
        }
    };

    const clearNotifications = () => {
        setNotificationCount(0);
    };

    const handleSearchFocus = () => {
        setIsSearching(true);
        setShowNotifications(false);
        setShowChat(false);
    };

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchQuery(e.target.value);
        if (e.target.value) {
            setIsSearching(true);
        }
    };

    const handleSubredditSelect = (subreddit: string) => {
        router.push(`/r/${subreddit}`);
        setIsSearching(false);
        setSearchQuery('');
    };

    const handleUserSelect = (username: string) => {
        router.push(`/user/${username}`);
        setIsSearching(false);
        setSearchQuery('');
    };

    return (
        <header className="bg-background sticky top-0 z-50 flex w-full items-center border-b">
            <div className="flex h-(--header-height) w-full items-center gap-2 px-4">
                <Button
                    className="h-8 w-8"
                    variant="ghost"
                    size="icon"
                    onClick={toggleSidebar}
                >
                    <SidebarIcon />
                </Button>
                {/* Logo */}
                <Link href="#" className="flex items-center gap-2">
                    <span className="text-xl font-bold hidden sm:inline">
                        VoxPopuli
                    </span>
                </Link>

                {/* Search */}
                <div className="relative flex-1 w-full ml-3">
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
                        <Input
                            ref={searchInputRef}
                            placeholder="Search VoxPopuli"
                            className="pl-9 bg-muted border-none focus-visible:ring-1 w-full"
                            value={searchQuery}
                            onChange={handleSearchChange}
                            onFocus={handleSearchFocus}
                        />
                        {searchQuery && (
                            <Button
                                variant="ghost"
                                size="icon"
                                className="absolute right-1 top-1/2 -translate-y-1/2 h-6 w-6"
                                onClick={() => {
                                    setSearchQuery('');
                                    setIsSearching(false);
                                }}
                            >
                                <X className="h-4 w-4" />
                            </Button>
                        )}
                    </div>

                    {isSearching && (
                        <div
                            ref={searchRef}
                            className="absolute top-full left-0 right-0 mt-1 bg-background border rounded-md shadow-lg z-50 max-h-[70vh] overflow-auto"
                            style={{
                                minHeight: searchQuery ? '200px' : '100px'
                            }}
                        >
                            {searchQuery ? (
                                <div className="p-2">
                                    <div className="mb-4">
                                        <h3 className="text-sm font-medium mb-2 px-2">
                                            Communities
                                        </h3>
                                        {searchResults.subreddits.length > 0 ? (
                                            <div className="space-y-1">
                                                {searchResults.subreddits.map(
                                                    (subreddit) => (
                                                        <div
                                                            key={subreddit.name}
                                                            className="flex items-center gap-2 p-2 hover:bg-muted rounded-md cursor-pointer"
                                                            onClick={() =>
                                                                handleSubredditSelect(
                                                                    subreddit.name
                                                                )
                                                            }
                                                        >
                                                            <Avatar className="h-6 w-6 flex-shrink-0">
                                                                <AvatarImage
                                                                    src={
                                                                        subreddit.icon
                                                                    }
                                                                    alt={`r/${subreddit.name}`}
                                                                />
                                                                <AvatarFallback>
                                                                    r/
                                                                </AvatarFallback>
                                                            </Avatar>
                                                            <div className="min-w-0">
                                                                <p className="text-sm font-medium truncate">
                                                                    r/
                                                                    {
                                                                        subreddit.name
                                                                    }
                                                                </p>
                                                                <p className="text-xs text-muted-foreground">
                                                                    {
                                                                        subreddit.members
                                                                    }
                                                                    M members
                                                                </p>
                                                            </div>
                                                        </div>
                                                    )
                                                )}
                                            </div>
                                        ) : (
                                            <p className="text-sm text-muted-foreground px-2">
                                                No communities found
                                            </p>
                                        )}
                                    </div>

                                    <div>
                                        <h3 className="text-sm font-medium mb-2 px-2">
                                            Users
                                        </h3>
                                        {searchResults.users.length > 0 ? (
                                            <div className="space-y-1">
                                                {searchResults.users.map(
                                                    (user) => (
                                                        <div
                                                            key={user.username}
                                                            className="flex items-center gap-2 p-2 hover:bg-muted rounded-md cursor-pointer"
                                                            onClick={() =>
                                                                handleUserSelect(
                                                                    user.username
                                                                )
                                                            }
                                                        >
                                                            <Avatar className="h-6 w-6 flex-shrink-0">
                                                                <AvatarImage
                                                                    src={
                                                                        user.avatar
                                                                    }
                                                                    alt={`u/${user.username}`}
                                                                />
                                                                <AvatarFallback>
                                                                    u/
                                                                </AvatarFallback>
                                                            </Avatar>
                                                            <div className="min-w-0">
                                                                <p className="text-sm font-medium truncate">
                                                                    u/
                                                                    {
                                                                        user.username
                                                                    }
                                                                </p>
                                                                <p className="text-xs text-muted-foreground">
                                                                    {user.karma.toLocaleString()}{' '}
                                                                    karma
                                                                </p>
                                                            </div>
                                                        </div>
                                                    )
                                                )}
                                            </div>
                                        ) : (
                                            <p className="text-sm text-muted-foreground px-2">
                                                No users found
                                            </p>
                                        )}
                                    </div>
                                </div>
                            ) : (
                                <div className="p-4 text-center text-muted-foreground">
                                    <p>
                                        Type to search for communities or users
                                    </p>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </header>
    );
}
