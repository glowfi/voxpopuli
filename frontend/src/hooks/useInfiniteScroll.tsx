'use client';
import React, { useState } from 'react';

export function useInfiniteScroll<T>(
    queryfn: (skip: number, limit: number) => Promise<T[]>,
    offset: number
) {
    const [page, setPage] = React.useState(0);
    const [loading, setLoading] = React.useState(false);
    const [hasMore, setHasMore] = React.useState(true);
    const [data, setData] = useState<T[]>([]);

    const next = async () => {
        if (loading) return;
        setLoading(true);

        try {
            const newData = await queryfn(page * offset, offset);

            setData((olddata) => [...olddata, ...newData]);
            setPage((prev) => prev + 1);

            if (newData.length < offset) {
                setHasMore(false);
            }
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    return { data, loading, hasMore, next };
}
