import { isLocal, proxyUrl } from './urls';
import * as cache from "./cache";
import { isStringInArray } from './arrays';


const proxyRequestInterceptor = (request) => {
    // Do something before request is sent
    if (!isLocal(request.url)) {
        request.url = proxyUrl(request.url);
        request.withCredentials = true;
    }
    return request;
}

function cacheResponseInterceptor(response, expires=null, whiteList=[]) {
    if (response.config.method === 'GET' || 'get') {
        if (response.config.url && !isStringInArray(response.config.url, whiteList)) {
            console.log(`caching: ${response.config.url}`);
            cache.store(response.config.url, response.data, expires);
        }
    }
    return response;
}

function cacheBlobResponseInterceptor(response, expires=null, whiteList=[]) {
    if (response.config.method === 'GET' || 'get') {
        if (response.config.url && !isStringInArray(response.config.url, whiteList)) {
            console.log(`caching: ${response.config.url}`);
            var objectUrl = URL.createObjectURL(response.data);
            cache.store(response.config.url, objectUrl, expires, 'blob');
        }
    }
    return response;
}

function cacheErrorInterceptor(error) {
    if (error.headers && error.headers.cached === true) {
        return Promise.resolve(error);
    }
    return Promise.reject(error);
}

function cacheRequestInterceptor(request, bucket='') {
    if (request.method === 'GET' || 'get') {
        const checkIsValidResponse = cache.isValid(request.url, bucket);
        if (checkIsValidResponse.isValid) {
            console.log(`cache hit: ${request.url}`);
            request.headers.cached = true;
            request.data = checkIsValidResponse.value;
            return Promise.reject(request);
        }
    }
    return request;
}

export {
    proxyRequestInterceptor,
    cacheResponseInterceptor,
    cacheRequestInterceptor,
    cacheErrorInterceptor,
    cacheBlobResponseInterceptor
};