import { TRPCError } from '@trpc/server';
import { verify, type JwtPayload } from 'jsonwebtoken';
import { t } from './trpc';
import type { Role } from './types';
import { authorizeUser } from './utils';

export const isAuthed = (...roles: Role[]) =>
    t.middleware(async (opts) => {
        // console.log(roles, 'Roles');
        // const { token } = opts.ctx;
        // if (!token) {
        //     throw new TRPCError({
        //         code: 'FORBIDDEN',
        //         message: 'Token not found.'
        //     });
        // }

        // let uid;

        // try {
        //     const user = verify(token, process.env.NEXTAUTH_SECRET || '');
        //     uid = (user as JwtPayload).uid;
        // } catch (error) {
        //     throw new TRPCError({
        //         code: 'FORBIDDEN',
        //         message: 'Invalid token.'
        //     });
        // }

        // await authorizeUser(uid, roles);

        // return opts.next({ ...opts, ctx: { ...opts.ctx, uid } });
        return opts.next({ ...opts, ctx: { ...opts.ctx } });
    });
