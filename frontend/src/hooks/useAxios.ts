import axios, { AxiosError, AxiosResponse } from 'axios';
import { useEffect, useState } from 'react';

axios.defaults.baseURL = process.env.API_ENDPOINT;

type Method = 'GET' | 'POST' | 'HEAD' | 'PUT' | 'PATCH' | 'DELETE';
type ResponseType = 'arraybuffer' | 'document' | 'json' | 'text' | 'stream';

export interface UseAxiosReturns {
    data: AxiosResponse | null;
    status: number | null;
    error: AxiosError | null;
    loading: boolean;
}

export interface UseAxiosArgs {
    url: string;
    method: Method;
    baseURL?: string;
    headers?: { [key: string]: string };
    params?: { [key: string]: string | number | boolean };
    payload?: string | number | boolean | object | null;
    timeout?: number;
    withCredentials?: boolean;
    responseType?: ResponseType;
}

const useAxios = ({
    url,
    method,
    baseURL,
    headers,
    params,
    payload,
    timeout,
    withCredentials,
    responseType
}: UseAxiosArgs): UseAxiosReturns => {
    const [data, setData] = useState<AxiosResponse | null>(null);
    const [status, setStatus] = useState<number | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<AxiosError | null>(null);

    useEffect(() => {
        axios({
            url,
            method,
            baseURL,
            headers,
            params,
            data: payload,
            timeout,
            withCredentials,
            responseType
        })
            .then(({ data, status }: AxiosResponse) => {
                setData(data);
                setStatus(status);
            })
            .catch((error: AxiosError) => setError(error))
            .finally(() => setLoading(false));
    }, [
        url,
        method,
        baseURL,
        headers,
        params,
        payload,
        timeout,
        withCredentials,
        responseType
    ]);

    return { data, status, error, loading };
};

export default useAxios;
