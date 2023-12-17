import request from '@/utils/request'

// get Vpn info
export function getInfo(params) {
    return request({
        url: '/api/vm/info',
        method: 'get',
        params
    })
}
export function ovpnRestart() {
    return request({
        url: `/api/ovpn/restart`,
        method: 'get',
    })
}
export function ovpnStatus() {
    return request({
        url: `/api/ovpn/status`,
        method: 'get',
    })
}
