import type { AppRouter } from '@voxpopuli/trpc-server/routers';
import { createTRPCReact } from '@trpc/react-query';

export const trpcClient = createTRPCReact<AppRouter>();
