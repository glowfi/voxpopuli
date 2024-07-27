import { PrismaClient } from '@prisma/client';
import 'dotenv/config';

const globalForPrisma = globalThis as unknown as {
    prisma: PrismaClient | undefined;
};

console.log(process.env.NODE_ENV, 'DB');

export const prisma =
    globalForPrisma.prisma ??
    new PrismaClient({
        log:
            process.env.NODE_ENV === 'development'
                ? ['query', 'error', 'warn']
                : ['error']
    });

if (process.env.NODE_ENV !== 'production') globalForPrisma.prisma = prisma;
