import { prisma } from '@voxpopuli/db';
import { z } from 'zod';
import { publicProcedure, router } from '../trpc';

export const topicsRouter = router({
    getAllTopics: publicProcedure
        .input(z.object({ skip: z.number(), limit: z.number() }))
        .query(async (req) => {
            let allTopics = await prisma.topics.findMany({
                select: {
                    title: true
                },
                skip: req.input.skip,
                take: req.input.limit
            });
            return allTopics;
        })
});
