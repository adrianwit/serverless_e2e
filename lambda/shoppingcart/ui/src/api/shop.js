const axios = require('axios');

let restClient;

function getRestClient() {
    if (restClient) {
        return restClient
    }
    restClient = axios.create({
        baseURL: document.URL,
        timeout: 1000
    });
    return restClient;
}




export default {
    getProducts: (cb) => {
        let client = getRestClient();
        client.get('/v1/api/products').then(function (resp) {
            cb(resp.data)
        }).catch(function (error) {
            cb({"Status":"error", "Error":error})
        });
    },
    
    buyProducts: (payload, cb, timeout) => setTimeout(() => cb(), timeout || 100)
}
