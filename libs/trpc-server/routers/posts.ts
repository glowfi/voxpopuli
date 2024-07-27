import { prisma } from '@voxpopuli/db';
import { z } from 'zod';
import { publicProcedure, router } from '../trpc';

export const postsRouter = router({
    getallposts: publicProcedure
        .input(z.object({ limit: z.number(), skip: z.number() }))
        .query(async (req) => {
            const { input } = req;
            let allPosts = await prisma.posts.findMany({
                include: {
                    author: {
                        select: {
                            username: true,
                            avatarImg: true
                        }
                    },
                    awards: {
                        select: {
                            _count: true,
                            title: true,
                            imageLink: true
                        }
                    },
                    mediaContent: {
                        select: {
                            postId: true,
                            id: true,
                            type: true,
                            gallery: true,
                            gifLink: true,
                            imageLink: true,
                            videoLink: true,
                            original_link: true
                        }
                    },
                    voxsphere: {
                        select: {
                            name: true,
                            logoURL: true
                        }
                    },
                    postflair: {
                        select: {
                            title: true,
                            colorHex: true
                        }
                    }
                },
                take: input.limit,
                skip: input.skip
            });
            return allPosts;
        })
});
