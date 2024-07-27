'use client';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import { LogOutIcon, RssIcon, SettingsIcon, UserIcon } from 'lucide-react';
import Link from 'next/link';
import { Button } from '../ui/button';
import SearchComplete from './search-complete';
import { SideBar } from './side-bar';
import { ModeToggle } from '../toggletheme/ThemeSwitcher';
import AppIcon from '@/app/icon.svg';
import Image from 'next/image';

const Navbar = () => {
    return (
        <header className="sticky top-0 z-50 flex items-center justify-between gap-4 bg-background px-4 py-2 shadow-sm md:px-6 border-b backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="flex items-center justify-between">
                <SideBar />
                <Link href="/" className="flex items-center gap-2">
                    <Image
                        className="hover:transition hover:opacity-75 hover:cursor-pointer"
                        src={AppIcon}
                        height={60}
                        width={60}
                        alt="Not Found"
                    />
                </Link>
            </div>
            <SearchComplete />
            <div className="flex justify-center items-center">
                <ModeToggle />
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="rounded-full"
                        >
                            <Avatar className="h-8 w-8">
                                <AvatarImage src="/placeholder-user.jpg" />
                                <AvatarFallback>JD</AvatarFallback>
                            </Avatar>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuLabel>
                            Signed in as John Doe
                        </DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem>
                            <UserIcon className="mr-2 h-4 w-4" />
                            Profile
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                            <SettingsIcon className="mr-2 h-4 w-4" />
                            Settings
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem>
                            <LogOutIcon className="mr-2 h-4 w-4" />
                            Logout
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </header>
    );
};

export default Navbar;
