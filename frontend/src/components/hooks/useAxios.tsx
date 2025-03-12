import axios, { AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { useState, useEffect } from 'react';

axios.defaults.baseURL = process.env.API_ENDPOINT;

interface UseAxiosOptions {
    url: string;
    method: 'get' | 'post' | 'put' | 'delete' | 'patch';
    body?: string;
    headers?: string;
    params?: { [key: string]: string | number | boolean }; // query parameters
    pathParams?: { [key: string]: string | number | boolean }; // path parameters
}

interface UseAxiosResponse {
    response: any;
    error: AxiosError | null;
    loading: boolean;
}

const useAxios = ({
    url,
    method,
    body,
    headers,
    params,
    pathParams
}: UseAxiosOptions): UseAxiosResponse => {
    const [response, setResponse] = useState<any>();
    const [error, setError] = useState<AxiosError | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    const fetchData = async () => {
        try {
            const config: AxiosRequestConfig = {
                method,
                url: getPathWithParams(url, pathParams),
                headers: headers ? JSON.parse(headers) : undefined,
                data: body ? JSON.parse(body) : undefined,
                params: params
            };

            const res: AxiosResponse = await axios(config);
            setResponse(res.data);
        } catch (err: any) {
            setError(err);
        } finally {
            setLoading(false);
        }
    };

    const getPathWithParams = (
        path: string,
        params: { [key: string]: string | number | boolean } | undefined
    ) => {
        if (!params) return path;
        return Object.keys(params).reduce(
            (acc, key) => acc.replace(`:${key}`, `${params[key]}`),
            path
        );
    };

    useEffect(() => {
        fetchData();
    }, [method, url, body, headers, params, pathParams]);

    return { response, error, loading };
};

export default useAxios;
