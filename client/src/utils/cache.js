import * as lscache from "lscache";


const CACHE_INTERVAL = 2 * 60 * 1000; // 2 minutes

function store(key, value, expires=CACHE_INTERVAL, bucket='') {
    if (bucket) {
        lscache.setBucket(bucket);
    }
    lscache.set(key, value, expires);
    if (bucket) {
        lscache.resetBucket();
    }
}

function isValid(key, bucket='') {
    if (bucket) {
        lscache.setBucket(bucket);
    }
    const value = lscache.get(key);
    if (bucket) {
        lscache.resetBucket();
    }
    return {
        isValid: value !== null,
        value: value,
    };
}

function initCache() {
    lscache.setExpiryMilliseconds(1);
    lscache.flushExpired();
    lscache.setBucket('blob');
    lscache.flush();
    lscache.resetBucket();
}

export {
    store,
    isValid,
    initCache
};