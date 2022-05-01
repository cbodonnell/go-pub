import axios from "axios";
import { cacheBlobResponseInterceptor, cacheErrorInterceptor, cacheRequestInterceptor, cacheResponseInterceptor, proxyRequestInterceptor } from "./interceptors";


class HttpClient {
    queue = {};

    client = axios.create();
    get(url, config=null) {
        if (!this.queue.hasOwnProperty(url)) {
            this.queue[url] = this.client.get(url, config).finally(() => delete this.queue[url]);
        } else {
            console.log(`request to ${url} already in progress`);
        }
        return this.queue[url];
    }
}

const proxyClient = new HttpClient();
proxyClient.client.interceptors.request.use(proxyRequestInterceptor);

const proxyCacheClient = new HttpClient();
proxyCacheClient.client.interceptors.request.use(
    (request) => {
        proxyRequestInterceptor(request);
        return cacheRequestInterceptor(request);
    },
);
proxyCacheClient.client.interceptors.response.use(
    (response) => cacheResponseInterceptor(response, 2*60*1000),
    (error) => cacheErrorInterceptor(error)
);

const proxyCacheBlobClient = new HttpClient();
proxyCacheBlobClient.client.interceptors.request.use(
    (request) => {
        proxyRequestInterceptor(request);
        return cacheRequestInterceptor(request, 'blob');
    },
);
proxyCacheBlobClient.client.interceptors.response.use(
    (response) => cacheBlobResponseInterceptor(response, 20*60*1000),
    (error) => cacheErrorInterceptor(error)
);



export {
    proxyClient,
    proxyCacheClient,
    proxyCacheBlobClient
};