import axiosInstance from '@/lib/axios-instance';

export interface PostAPI {
    posts(skip: number, limit: number): Promise<Post[]>;
}

class PostClient implements PostAPI {
    async posts(skip: number, limit: number): Promise<Post[]> {
        const resp = await axiosInstance()({
            url: '/api/posts',
            method: 'get',
            params: {
                skip: skip,
                limit: limit
            }
        });

        return resp.data as Post[];
    }
}

export const postApi = new PostClient();
