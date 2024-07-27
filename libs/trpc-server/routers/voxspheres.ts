import { prisma } from '@voxpopuli/db';
import { z } from 'zod';
import { publicProcedure, router } from '../trpc';

export const voxspheresRouter = router({
    getTopVoxspheres: publicProcedure
        .input(z.object({ limit: z.number() }))
        .query(async (req) => {
            let topVoxSpheres = await prisma.voxsphere.findMany({
                select: {
                    totalmembers: true,
                    totalmembersHuman: true,
                    name: true,
                    logoURL: true
                },
                orderBy: {
                    totalmembers: 'desc'
                },
                take: req.input.limit
            });

            return topVoxSpheres;
        })
});
