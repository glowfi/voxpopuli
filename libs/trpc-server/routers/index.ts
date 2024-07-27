import type { inferRouterOutputs } from '@trpc/server';
import { router } from '../trpc';
import { postsRouter } from './posts';
import { topicsRouter } from './topics';
import { voxspheresRouter } from './voxspheres';

export const appRouter = router({
    posts: postsRouter,
    voxspheres: voxspheresRouter,
    topics: topicsRouter
});

export type AppRouter = typeof appRouter;
export type AppRouterType = inferRouterOutputs<AppRouter>;
