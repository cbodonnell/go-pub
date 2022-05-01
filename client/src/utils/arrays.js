function isStringInArray(url, list) {
    return list.findIndex(w => url.includes(w)) > -1;
}

export {
    isStringInArray
}