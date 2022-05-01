function formDataToURLString(formData) {
    const data = [...formData.entries()];
    const asString = data
    .map(x => `${encodeURIComponent(x[0])}=${encodeURIComponent(x[1])}`)
    .join('&');
    return asString;
}

let lastId = 0;

function uniqueID(prefix='') {
    lastId++;
    return `${prefix}${lastId}`;
}

export {
    formDataToURLString,
    uniqueID
}