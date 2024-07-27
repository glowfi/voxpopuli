import { initTRPC } from '@trpc/server';
import type { TRPCContext } from './context';
import { isAuthed } from './middleware';
import type { Role } from './types';

export const t = initTRPC.context<TRPCContext>().create();

export const router = t.router;
export const publicProcedure = t.procedure;
// export const privateProcedure = t.procedure.use(isAuthed());
export const privateProcedure = t.procedure.use(isAuthed('admin', 'manager'));
// export const privateProcedure = (...roles: Role[]) =>
//     t.procedure.use(isAuthed(...roles));
