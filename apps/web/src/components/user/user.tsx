import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle
} from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { AwardIcon, ClubIcon, TrophyIcon } from 'lucide-react';

export default function User() {
    return (
        <div className="container mx-auto p-4 grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="md:col-span-2 space-y-4">
                <Tabs defaultValue="posts">
                    <TabsList className="grid grid-cols-3">
                        <TabsTrigger value="posts">Posts</TabsTrigger>
                        <TabsTrigger value="comments">Comments</TabsTrigger>
                        <TabsTrigger value="overview">Overview</TabsTrigger>
                    </TabsList>
                    <TabsContent value="posts">
                        <Card>
                            <CardHeader>
                                <CardTitle>Post Title 1</CardTitle>
                                <CardDescription>
                                    Posted by u/spez
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p>This is the content of the first post.</p>
                            </CardContent>
                            <CardFooter>
                                <Button variant="outline">Comment</Button>
                            </CardFooter>
                        </Card>
                        <Card>
                            <CardHeader>
                                <CardTitle>Post Title 2</CardTitle>
                                <CardDescription>
                                    Posted by u/spez
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p>This is the content of the second post.</p>
                            </CardContent>
                            <CardFooter>
                                <Button variant="outline">Comment</Button>
                            </CardFooter>
                        </Card>
                    </TabsContent>
                    <TabsContent value="comments">
                        <Card>
                            <CardHeader>
                                <CardTitle>Comment on Post Title 1</CardTitle>
                                <CardDescription>
                                    Commented by u/spez
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p>This is the content of the first comment.</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader>
                                <CardTitle>Comment on Post Title 2</CardTitle>
                                <CardDescription>
                                    Commented by u/spez
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p>
                                    This is the content of the second comment.
                                </p>
                            </CardContent>
                        </Card>
                    </TabsContent>
                    <TabsContent value="overview">
                        <Card>
                            <CardHeader>
                                <CardTitle>Overview</CardTitle>
                                <CardDescription>
                                    Summary of user activity
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p>
                                    This is an overview of the users activity.
                                </p>
                            </CardContent>
                        </Card>
                    </TabsContent>
                </Tabs>
            </div>
            <div className="space-y-4">
                <Card>
                    <CardHeader>
                        <CardTitle>User Details</CardTitle>
                    </CardHeader>
                    <CardContent className="flex items-center space-x-4">
                        <Avatar>
                            <AvatarImage
                                src="/images/placeholder.svg"
                                alt="User Avatar"
                                width={100}
                                height={100}
                            />
                            <AvatarFallback>SP</AvatarFallback>
                        </Avatar>
                        <div>
                            <h2 className="text-xl font-bold">u/spez</h2>
                            <p>Reddit Admin</p>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Trophies</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <div className="flex items-center space-x-2">
                            <TrophyIcon />
                            <p>Trophy 1</p>
                        </div>
                        <div className="flex items-center space-x-2">
                            <TrophyIcon />
                            <p>Trophy 2</p>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Basic Info</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <div>
                            <Label>Joined:</Label>
                            <p>January 1, 2006</p>
                        </div>
                        <div>
                            <Label>Karma:</Label>
                            <p>1,234,567</p>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Awards</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <div className="flex items-center space-x-2">
                            <AwardIcon />
                            <p>Award 1</p>
                        </div>
                        <div className="flex items-center space-x-2">
                            <AwardIcon />
                            <p>Award 2</p>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Moderated Communities</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <div className="flex items-center space-x-2">
                            <ClubIcon />
                            <p>r/community1</p>
                        </div>
                        <div className="flex items-center space-x-2">
                            <ClubIcon />
                            <p>r/community2</p>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
