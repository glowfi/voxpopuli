'use client';

import { useState, useEffect } from 'react';

const useInfiniteScroll = <T,>(
    queryFn: (skip: number, limit: number) => Promise<T[]>,
    offset: number
) => {
    const [data, setData] = useState<T[]>([]);
    const [loading, setLoading] = useState(false);
    const [hasNextPage, setHasNextPage] = useState(true);
    const [skip, setSkip] = useState(0);

    const fetchMoreData = async () => {
        if (loading || !hasNextPage) return;
        setLoading(true);
        try {
            // await new Promise((r) => setTimeout(r, 500));
            const newData = await queryFn(skip, offset);
            setData([...data, ...newData]);
            setHasNextPage(newData.length === offset);
            setSkip(skip + offset);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleScroll = () => {
        const scrollPosition = window.scrollY + window.innerHeight;
        const height = document.body.offsetHeight;
        if (scrollPosition >= height * 0.9) {
            fetchMoreData();
        }
    };

    useEffect(() => {
        window.addEventListener('scroll', handleScroll);
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, [hasNextPage, loading, skip, offset]);

    useEffect(() => {
        fetchMoreData();
    }, []);

    return { data, loading, hasNextPage };
};

export default useInfiniteScroll;
