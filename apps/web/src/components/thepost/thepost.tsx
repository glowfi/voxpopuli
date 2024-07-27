import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle
} from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import Image from 'next/image';
import Link from 'next/link';

export default function ThePost() {
    return (
        <div className="grid grid-cols-1 md:grid-cols-[3fr_1fr] gap-8 max-w-6xl mx-auto py-8 px-4">
            <div className="space-y-6">
                <div>
                    <h1 className="text-3xl font-bold">
                        The Joke Tax Chronicles
                    </h1>
                    <div className="flex items-center gap-2 text-muted-foreground">
                        <div className="flex items-center gap-1">
                            <Avatar className="w-5 h-5 border">
                                <AvatarImage src="/placeholder-user.jpg" />
                                <AvatarFallback>AC</AvatarFallback>
                            </Avatar>
                            <span>u/shadcn</span>
                        </div>
                        <Separator orientation="vertical" />
                        <span>Posted 2 days ago</span>
                    </div>
                </div>
                <div className="prose prose-gray dark:prose-invert">
                    <p>
                        Once upon a time, in a far-off land, there was a very
                        lazy king who spent all day lounging on his throne. One
                        day, his advisors came to him with a problem: the
                        kingdom was running out of money.
                    </p>
                    <p>
                        Jokester began sneaking into the castle in the middle of
                        the night and leaving jokes all over the place: under
                        the king&apos;s pillow, in his soup, even in the royal
                        toilet. The king was furious, but he couldn&apos;t seem
                        to stop Jokester.
                    </p>
                    <p>
                        And then, one day, the people of the kingdom discovered
                        that the jokes left by Jokester were so funny that they
                        couldn&apos;t help but laugh. And once they started
                        laughing, they couldn&apos;t stop.
                    </p>
                    <blockquote>
                        &ldquo;After all,&rdquo; he said, &ldquo;everyone enjoys
                        a good joke, so it&apos;s only fair that they should pay
                        for the privilege.&rdquo;
                    </blockquote>
                    <h3>The Joke Tax</h3>
                    <p>
                        The king&apos;s subjects were not amused. They grumbled
                        and complained, but the king was firm:
                    </p>
                    <ul>
                        <li>1st level of puns: 5 gold coins</li>
                        <li>2nd level of jokes: 10 gold coins</li>
                        <li>3rd level of one-liners : 20 gold coins</li>
                    </ul>
                    <p>
                        As a result, people stopped telling jokes, and the
                        kingdom fell into a gloom. But there was one person who
                        refused to let the king&apos;s foolishness get him down:
                        a court jester named Jokester.
                    </p>
                </div>
                <div className="space-y-4">
                    <h2 className="text-2xl font-bold">Comments</h2>
                    <div className="space-y-4">
                        <div className="flex items-start gap-4">
                            <Avatar className="w-8 h-8 border">
                                <AvatarImage src="/placeholder-user.jpg" />
                                <AvatarFallback>AC</AvatarFallback>
                            </Avatar>
                            <div className="space-y-2">
                                <div className="flex items-center gap-2">
                                    <div className="font-medium">u/shadcn</div>
                                    <span className="text-muted-foreground">
                                        2 days ago
                                    </span>
                                </div>
                                <p>
                                    This is a hilarious story! I cant believe
                                    the king tried to tax jokes. What a
                                    ridiculous idea.
                                </p>
                                <div className="flex items-center gap-2">
                                    <Button variant="ghost" size="icon">
                                        <ThumbsUpIcon className="w-4 h-4" />
                                        <span className="sr-only">Upvote</span>
                                    </Button>
                                    <Button variant="ghost" size="icon">
                                        <ThumbsDownIcon className="w-4 h-4" />
                                        <span className="sr-only">
                                            Downvote
                                        </span>
                                    </Button>
                                    <Button variant="ghost" size="icon">
                                        <MessageCircleIcon className="w-4 h-4" />
                                        <span className="sr-only">Reply</span>
                                    </Button>
                                </div>
                            </div>
                        </div>
                        <div className="flex items-start gap-4 pl-12">
                            <Avatar className="w-8 h-8 border">
                                <AvatarImage src="/placeholder-user.jpg" />
                                <AvatarFallback>AM</AvatarFallback>
                            </Avatar>
                            <div className="space-y-2">
                                <div className="flex items-center gap-2">
                                    <div className="font-medium">u/amelia</div>
                                    <span className="text-muted-foreground">
                                        1 day ago
                                    </span>
                                </div>
                                <p>
                                    I cant believe the king was so out of touch
                                    with his people. Taxing jokes? Thats just
                                    cruel!
                                </p>
                                <div className="flex items-center gap-2">
                                    <Button variant="ghost" size="icon">
                                        <ThumbsUpIcon className="w-4 h-4" />
                                        <span className="sr-only">Upvote</span>
                                    </Button>
                                    <Button variant="ghost" size="icon">
                                        <ThumbsDownIcon className="w-4 h-4" />
                                        <span className="sr-only">
                                            Downvote
                                        </span>
                                    </Button>
                                    <Button variant="ghost" size="icon">
                                        <MessageCircleIcon className="w-4 h-4" />
                                        <span className="sr-only">Reply</span>
                                    </Button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div className="space-y-6">
                <Card>
                    <CardHeader>
                        <CardTitle>r/Jokes</CardTitle>
                        <CardDescription>
                            A subreddit dedicated to sharing and discussing
                            jokes.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2">
                            <div className="flex items-center gap-2">
                                <Avatar className="w-8 h-8 border">
                                    <AvatarImage src="/placeholder-user.jpg" />
                                    <AvatarFallback>J</AvatarFallback>
                                </Avatar>
                                <div>
                                    <div className="font-medium">r/Jokes</div>
                                    <div className="text-muted-foreground">
                                        A community for sharing and discussing
                                        jokes.
                                    </div>
                                </div>
                            </div>
                            <div className="grid gap-2">
                                <Link
                                    href="#"
                                    className="flex items-center gap-2 text-sm"
                                    prefetch={false}
                                >
                                    <ThumbsUpIcon className="w-4 h-4" />
                                    <span>Top Jokes</span>
                                </Link>
                                <Link
                                    href="#"
                                    className="flex items-center gap-2 text-sm"
                                    prefetch={false}
                                >
                                    <TrendingUpIcon className="w-4 h-4" />
                                    <span>Trending Jokes</span>
                                </Link>
                                <Link
                                    href="#"
                                    className="flex items-center gap-2 text-sm"
                                    prefetch={false}
                                >
                                    <BookmarkIcon className="w-4 h-4" />
                                    <span>Saved Jokes</span>
                                </Link>
                            </div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Related Content</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid gap-2">
                            <Link
                                href="#"
                                className="flex items-center gap-2 text-sm"
                                prefetch={false}
                            >
                                <Image
                                    src="/placeholder.svg"
                                    alt="Related post"
                                    width={48}
                                    height={48}
                                    className="rounded-md"
                                />
                                <div>
                                    <div className="font-medium">
                                        The Funniest Jokes of 2023
                                    </div>
                                    <div className="text-muted-foreground">
                                        r/Jokes
                                    </div>
                                </div>
                            </Link>
                            <Link
                                href="#"
                                className="flex items-center gap-2 text-sm"
                                prefetch={false}
                            >
                                <Image
                                    src="/placeholder.svg"
                                    alt="Related post"
                                    width={48}
                                    height={48}
                                    className="rounded-md"
                                />
                                <div>
                                    <div className="font-medium">
                                        How to Tell a Joke: A Beginners Guide
                                    </div>
                                    <div className="text-muted-foreground">
                                        r/Jokes
                                    </div>
                                </div>
                            </Link>
                            <Link
                                href="#"
                                className="flex items-center gap-2 text-sm"
                                prefetch={false}
                            >
                                <Image
                                    src="/placeholder.svg"
                                    alt="Related post"
                                    width={48}
                                    height={48}
                                    className="rounded-md"
                                />
                                <div>
                                    <div className="font-medium">
                                        The Science Behind Laughter: Why Jokes
                                        are Funny
                                    </div>
                                    <div className="text-muted-foreground">
                                        r/Jokes
                                    </div>
                                </div>
                            </Link>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}

function BookmarkIcon(props: any) {
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
            <path d="m19 21-7-4-7 4V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z" />
        </svg>
    );
}

function MessageCircleIcon(props: any) {
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
            <path d="M7.9 20A9 9 0 1 0 4 16.1L2 22Z" />
        </svg>
    );
}

function ThumbsDownIcon(props: any) {
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
            <path d="M17 14V2" />
            <path d="M9 18.12 10 14H4.17a2 2 0 0 1-1.92-2.56l2.33-8A2 2 0 0 1 6.5 2H20a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2h-2.76a2 2 0 0 0-1.79 1.11L12 22h0a3.13 3.13 0 0 1-3-3.88Z" />
        </svg>
    );
}

function ThumbsUpIcon(props: any) {
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
            <path d="M7 10v12" />
            <path d="M15 5.88 14 10h5.83a2 2 0 0 1 1.92 2.56l-2.33 8A2 2 0 0 1 17.5 22H4a2 2 0 0 1-2-2v-8a2 2 0 0 1 2-2h2.76a2 2 0 0 0 1.79-1.11L12 2h0a3.13 3.13 0 0 1 3 3.88Z" />
        </svg>
    );
}

function TrendingUpIcon(props: any) {
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
            <polyline points="22 7 13.5 15.5 8.5 10.5 2 17" />
            <polyline points="16 7 22 7 22 13" />
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
