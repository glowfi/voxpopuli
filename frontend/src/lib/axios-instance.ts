import axios, { AxiosInstance, AxiosRequestConfig } from "axios";

export default function axiosInstance(
  config: AxiosRequestConfig = {},
): AxiosInstance {
  axios.defaults.baseURL = process.env.API_BASE_URL;

  return axios.create(config);
}
