import request from '@/utils/request'

// get Vpn info
export function getInfo(params) {
    return request({
        url: '/api/vm/info',
        method: 'get',
        params
    })
}
