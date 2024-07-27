import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import Link from 'next/link';
import { Separator } from '@/components/ui/separator';
import Image from 'next/image';

export default function Voxsphere() {
    return (
        <div className="flex flex-col min-h-dvh">
            <header className="bg-primary text-primary-foreground py-8 md:py-12">
                <div className="container px-4 md:px-6">
                    <div className="grid gap-4 md:gap-6 lg:gap-8">
                        <div>
                            <h1 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl lg:text-6xl">
                                v/science
                            </h1>
                            <p className="max-w-[720px] text-lg md:text-xl lg:text-2xl">
                                A community for sharing and discussing new and
                                high-quality science, technology, and health
                                information.
                            </p>
                        </div>
                        <div className="flex items-center gap-4 md:gap-6">
                            <div className="flex items-center gap-2">
                                <UsersIcon className="w-5 h-5" />
                                <span className="text-sm md:text-base">
                                    2.5M Members
                                </span>
                            </div>
                            <div className="flex items-center gap-2">
                                <ActivityIcon className="w-5 h-5" />
                                <span className="text-sm md:text-base">
                                    Active Now
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </header>
            <div className="container grid grid-cols-[240px_1fr] gap-8 px-4 py-8 md:px-6 md:py-12">
                <div className="space-y-6">
                    <Card>
                        <CardHeader>
                            <CardTitle>About v/science</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p>
                                r/science is a community for sharing and
                                discussing new and high-quality science,
                                technology, and health information.
                            </p>
                            <div className="mt-4 flex items-center gap-2">
                                <UsersIcon className="w-5 h-5" />
                                <span>2.5M Members</span>
                            </div>
                            <div className="mt-2 flex items-center gap-2">
                                <ActivityIcon className="w-5 h-5" />
                                <span>Active Now</span>
                            </div>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardHeader>
                            <CardTitle>Moderators</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="grid gap-4">
                                <div className="flex items-center gap-2">
                                    <Avatar className="w-8 h-8">
                                        <AvatarImage src="/placeholder-user.jpg" />
                                        <AvatarFallback>AC</AvatarFallback>
                                    </Avatar>
                                    <div>
                                        <div className="font-medium">
                                            u/shadcn
                                        </div>
                                        <div className="text-xs text-muted-foreground">
                                            Head Moderator
                                        </div>
                                    </div>
                                </div>
                                <div className="flex items-center gap-2">
                                    <Avatar className="w-8 h-8">
                                        <AvatarImage src="/placeholder-user.jpg" />
                                        <AvatarFallback>AC</AvatarFallback>
                                    </Avatar>
                                    <div>
                                        <div className="font-medium">
                                            u/acme
                                        </div>
                                        <div className="text-xs text-muted-foreground">
                                            Moderator
                                        </div>
                                    </div>
                                </div>
                                <div className="flex items-center gap-2">
                                    <Avatar className="w-8 h-8">
                                        <AvatarImage src="/placeholder-user.jpg" />
                                        <AvatarFallback>AC</AvatarFallback>
                                    </Avatar>
                                    <div>
                                        <div className="font-medium">
                                            u/vercel
                                        </div>
                                        <div className="text-xs text-muted-foreground">
                                            Moderator
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardHeader>
                            <CardTitle>Related Communities</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="grid gap-2">
                                <Link
                                    href="#"
                                    className="flex items-center gap-2 hover:underline"
                                    prefetch={false}
                                >
                                    <RssIcon className="w-5 h-5" />
                                    <span>r/technology</span>
                                </Link>
                                <Link
                                    href="#"
                                    className="flex items-center gap-2 hover:underline"
                                    prefetch={false}
                                >
                                    <RssIcon className="w-5 h-5" />
                                    <span>r/futurology</span>
                                </Link>
                                <Link
                                    href="#"
                                    className="flex items-center gap-2 hover:underline"
                                    prefetch={false}
                                >
                                    <RssIcon className="w-5 h-5" />
                                    <span>r/space</span>
                                </Link>
                                <Link
                                    href="#"
                                    className="flex items-center gap-2 hover:underline"
                                    prefetch={false}
                                >
                                    <RssIcon className="w-5 h-5" />
                                    <span>r/environment</span>
                                </Link>
                            </div>
                        </CardContent>
                    </Card>
                </div>
                <div>
                    <div className="grid gap-6 md:gap-8">
                        <Card>
                            <CardContent>
                                <div className="grid grid-cols-[100px_1fr] gap-4">
                                    <Image
                                        src="/placeholder.svg"
                                        alt="Post thumbnail"
                                        width={100}
                                        height={75}
                                        className="rounded-md object-cover"
                                    />
                                    <div className="space-y-2">
                                        <h3 className="text-lg font-medium">
                                            New study shows the benefits of
                                            meditation for mental health
                                        </h3>
                                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                            <div>u/acme</div>
                                            <Separator
                                                orientation="vertical"
                                                className="h-4"
                                            />
                                            <div>12 Comments</div>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent>
                                <div className="grid grid-cols-[100px_1fr] gap-4">
                                    <Image
                                        src="/placeholder.svg"
                                        alt="Post thumbnail"
                                        width={100}
                                        height={75}
                                        className="rounded-md object-cover"
                                    />
                                    <div className="space-y-2">
                                        <h3 className="text-lg font-medium">
                                            Breakthrough in fusion energy
                                            production
                                        </h3>
                                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                            <div>u/shadcn</div>
                                            <Separator
                                                orientation="vertical"
                                                className="h-4"
                                            />
                                            <div>42 Comments</div>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent>
                                <div className="grid grid-cols-[100px_1fr] gap-4">
                                    <Image
                                        src="/placeholder.svg"
                                        alt="Post thumbnail"
                                        width={100}
                                        height={75}
                                        className="rounded-md object-cover"
                                    />
                                    <div className="space-y-2">
                                        <h3 className="text-lg font-medium">
                                            NASAs Perseverance rover makes new
                                            discoveries on Mars
                                        </h3>
                                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                            <div>u/vercel</div>
                                            <Separator
                                                orientation="vertical"
                                                className="h-4"
                                            />
                                            <div>28 Comments</div>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent>
                                <div className="grid grid-cols-[100px_1fr] gap-4">
                                    <Image
                                        src="/placeholder.svg"
                                        alt="Post thumbnail"
                                        width={100}
                                        height={75}
                                        className="rounded-md object-cover"
                                    />
                                    <div className="space-y-2">
                                        <h3 className="text-lg font-medium">
                                            New study links gut microbiome to
                                            mental health
                                        </h3>
                                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                            <div>u/acme</div>
                                            <Separator
                                                orientation="vertical"
                                                className="h-4"
                                            />
                                            <div>19 Comments</div>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent>
                                <div className="grid grid-cols-[100px_1fr] gap-4">
                                    <Image
                                        src="/placeholder.svg"
                                        alt="Post thumbnail"
                                        width={100}
                                        height={75}
                                        className="rounded-md object-cover"
                                    />
                                    <div className="space-y-2">
                                        <h3 className="text-lg font-medium">
                                            Advances in renewable energy
                                            technology
                                        </h3>
                                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                            <div>u/shadcn</div>
                                            <Separator
                                                orientation="vertical"
                                                className="h-4"
                                            />
                                            <div>36 Comments</div>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent>
                                <div className="grid grid-cols-[100px_1fr] gap-4">
                                    <Image
                                        src="/placeholder.svg"
                                        alt="Post thumbnail"
                                        width={100}
                                        height={75}
                                        className="rounded-md object-cover"
                                    />
                                    <div className="space-y-2">
                                        <h3 className="text-lg font-medium">
                                            New findings on the impact of
                                            climate change
                                        </h3>
                                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                            <div>u/vercel</div>
                                            <Separator
                                                orientation="vertical"
                                                className="h-4"
                                            />
                                            <div>54 Comments</div>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </div>
            </div>
        </div>
    );
}

function ActivityIcon(props: any) {
    return (
        <svg
            {...props}
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <path d="M22 12h-2.48a2 2 0 0 0-1.93 1.46l-2.35 8.36a.25.25 0 0 1-.48 0L9.24 2.18a.25.25 0 0 0-.48 0l-2.35 8.36A2 2 0 0 1 4.49 12H2" />
        </svg>
    );
}

function RssIcon(props: any) {
    return (
        <svg
            {...props}
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <path d="M4 11a9 9 0 0 1 9 9" />
            <path d="M4 4a16 16 0 0 1 16 16" />
            <circle cx="5" cy="19" r="1" />
        </svg>
    );
}

function UsersIcon(props: any) {
    return (
        <svg
            {...props}
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
            <circle cx="9" cy="7" r="4" />
            <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
        </svg>
    );
}

function XIcon(props: any) {
    return (
        <svg
            {...props}
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
        </svg>
    );
}
