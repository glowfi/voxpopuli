'use client';
import useAxios from '@/components/hooks/useAxios';
import React, { useEffect, useState } from 'react';

const Posts = () => {
    const { response, loading, error } = useAxios({
        method: 'get',
        url: '/posts?skip=10&limit=10'
    });
    const [data, setData] = useState([]);

    useEffect(() => {
        if (response !== null) {
            setData(response);
        }
    }, [response]);

    return (
        <div className="App">
            <h1>Posts</h1>

            {loading ? (
                <p>loading...</p>
            ) : (
                <div>
                    {error && (
                        <div>
                            <p>{error.message}</p>
                        </div>
                    )}
                    <div>{data && <p>{JSON.stringify(data)}</p>}</div>
                </div>
            )}
        </div>
    );
};

export default Posts;
