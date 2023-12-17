import request from '@/utils/request'

export function getUsers(params) {
    return request({
        url: '/api/user/show',
        method: 'get',
        timeout: 60000,
        params
    })
}

export function createUser(data) {
    return request({
        url: '/api/user/add',
        method: 'post',
        data
    })
}

export function deleteUser(username) {
    return request({
        url: `/api/user/revoke/${username}`,
        method: 'delete',
    })
}