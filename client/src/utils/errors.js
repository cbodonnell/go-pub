export default function logError(error) {
    if (error.response) {
        // Request made and server responded
        console.error(`Error code ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
        // The request was made but no response was received
        console.error(error.request);
    } else {
        // Something happened in setting up the request that triggered an Error
        console.error('Error', error.message);
    }
}