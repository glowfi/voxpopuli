import React from 'react';
import { TOTAL_POSTS_TO_FETCH } from '../constants';
import { Skeleton } from '@/components/ui/skeleton';
import { Card } from '@/components/ui/card';
import { Separator } from '@radix-ui/react-dropdown-menu';

const PostSkeleton = () => {
    return Array.from({ length: TOTAL_POSTS_TO_FETCH }, (_, i) => i + 1).map(
        (p, idx) => {
            return (
                <Card className="flex items-start gap-4 p-4" key={idx}>
                    <Skeleton className="h-14 w-14 md:h-24 md:w-24 rounded-md" />
                    <div className="space-y-2 w-full">
                        <Skeleton className="h-4 w-[100px] md:w-[200px]" />
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <Skeleton className="h-4 w-[50px] md:w-[100px]" />
                            <Separator className="h-4" />
                            <Skeleton className="h-4 w-[40px] md:w-[80px]" />
                        </div>
                        <Skeleton className="h-4 w-[150px] md:w-[300px]" />
                        <Skeleton className="h-4 w-[125px] md:w-[250px]" />
                    </div>
                </Card>
            );
        }
    );
};

export default PostSkeleton;
