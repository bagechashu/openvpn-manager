import axios from 'axios';

// create an axios instance
const instance = axios.create({
    baseURL: process.env.NODE_ENV === 'production' ? process.env.VUE_APP_BASE_API : '/', // api base_url
    timeout: 5000 // request timeout
})

// request interceptor
instance.interceptors.request.use(
    config => {
        return config
    },
    error => {
        console.log(error) // for debug
        return Promise.reject(error)
    }
)

instance.interceptors.response.use(
    response => {
        const { status, data } = response
        switch (status) {
            case 209:
                window.$message.warning(data.response.msg, { duration: 2000 })
        }
        return Promise.resolve(data)
    },
    error => {
        const { status, data } = error.response
        switch (status) {
            case 401:
            case 403:
                window.$message.error('No Permission, Please Login!', { duration: 2000 })
                break
            default:
                window.$message.error(data.response.msg || error.message, { duration: 2000 })
        }
        return Promise.reject(error);
    }
)

export default instance