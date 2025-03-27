import { AppSidebar } from '@/components/app-sidebar';
import { SiteHeader } from '@/components/site-header';
import { Button } from '@/components/ui/button';
import {
    Card,
    CardHeader,
    CardTitle,
    CardContent,
    CardFooter,
    CardDescription
} from '@/components/ui/card';
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar';
import { ResponsivePostCard } from '@/post/ResposivePostCard';
import { ScrollArea } from '@radix-ui/react-scroll-area';
import {
    Home,
    Rocket,
    TrendingUp,
    ChevronUp,
    Plus,
    Sparkles
} from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';

export default function Page() {
    const communities = [
        {
            name: 'r/nextjs',
            members: '125.4k',
            icon: '/placeholder.svg?height=40&width=40'
        },
        {
            name: 'r/reactjs',
            members: '342.8k',
            icon: '/placeholder.svg?height=40&width=40'
        },
        {
            name: 'r/webdev',
            members: '1.2m',
            icon: '/placeholder.svg?height=40&width=40'
        },
        {
            name: 'r/programming',
            members: '5.3m',
            icon: '/placeholder.svg?height=40&width=40'
        },
        {
            name: 'r/tailwindcss',
            members: '87.2k',
            icon: '/placeholder.svg?height=40&width=40'
        }
    ];
    return (
        <div className="[--header-height:calc(--spacing(14))]">
            <SidebarProvider className="flex flex-col">
                <SiteHeader />
                <div className="flex flex-1">
                    <AppSidebar />
                    <SidebarInset>
                        <div className="container mx-auto grid grid-cols-1 md:grid-cols-12 m-3 p-3 gap-6 items-start justify-center">
                            <div className="col-span-1 md:col-span-8">
                                <ResponsivePostCard />
                                <ResponsivePostCard />
                                <ResponsivePostCard />
                                <ResponsivePostCard />
                                <ResponsivePostCard />
                            </div>
                            <div className="col-span-1 md:col-span-4">
                                <div className="space-y-4">
                                    <Card>
                                        <CardHeader className="pb-2">
                                            <CardTitle className="text-base">
                                                Top Communities
                                            </CardTitle>
                                        </CardHeader>
                                        <CardContent className="p-3 pt-0">
                                            <ScrollArea className="pr-4">
                                                <div className="space-y-1">
                                                    {communities.map(
                                                        (community, i) => (
                                                            <Link
                                                                key={
                                                                    community.name
                                                                }
                                                                href={`/${community.name}`}
                                                                className="flex items-center gap-3 rounded-md px-2 py-2 hover:bg-muted"
                                                            >
                                                                <span className="text-sm font-medium text-muted-foreground w-5">
                                                                    {i + 1}
                                                                </span>
                                                                <ChevronUp className="h-4 w-4 text-green-500" />
                                                                <div className="relative h-6 w-6 overflow-hidden rounded-full">
                                                                    <Image
                                                                        src={
                                                                            community.icon ||
                                                                            '/placeholder.svg'
                                                                        }
                                                                        alt={
                                                                            community.name
                                                                        }
                                                                        fill
                                                                        className="object-cover"
                                                                    />
                                                                </div>
                                                                <div className="flex-1 truncate">
                                                                    <div className="text-sm font-medium">
                                                                        {
                                                                            community.name
                                                                        }
                                                                    </div>
                                                                </div>
                                                            </Link>
                                                        )
                                                    )}
                                                </div>
                                            </ScrollArea>
                                        </CardContent>
                                        <CardFooter className="border-t p-3">
                                            <Button
                                                className="w-full"
                                                size="sm"
                                            >
                                                View All
                                            </Button>
                                        </CardFooter>
                                    </Card>

                                    {/* <Card> */}
                                    {/*     <CardHeader className="pb-2"> */}
                                    {/*         <CardTitle className="text-base"> */}
                                    {/*             Create Post */}
                                    {/*         </CardTitle> */}
                                    {/*         <CardDescription> */}
                                    {/*             Share your thoughts with the */}
                                    {/*             community */}
                                    {/*         </CardDescription> */}
                                    {/*     </CardHeader> */}
                                    {/*     <CardContent className="p-3 pt-0"> */}
                                    {/*         <Button className="w-full gap-2"> */}
                                    {/*             <Plus className="h-4 w-4" /> */}
                                    {/*             Create Post */}
                                    {/*         </Button> */}
                                    {/*     </CardContent> */}
                                    {/* </Card> */}

                                    <Card>
                                        <CardHeader className="pb-2">
                                            <CardTitle className="text-base">
                                                Premium
                                            </CardTitle>
                                            <CardDescription>
                                                Upgrade your VoxPopuli
                                                experience
                                            </CardDescription>
                                        </CardHeader>
                                        <CardContent className="p-3 pt-0">
                                            <Button
                                                variant="outline"
                                                className="w-full gap-2"
                                            >
                                                <Sparkles className="h-4 w-4 text-yellow-500" />
                                                Try Premium
                                            </Button>
                                        </CardContent>
                                    </Card>
                                </div>
                            </div>
                        </div>
                    </SidebarInset>
                </div>
            </SidebarProvider>
        </div>
    );
}
