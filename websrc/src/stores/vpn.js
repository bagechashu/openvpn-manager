import { defineStore } from 'pinia'
import { reactive } from 'vue'
import { getInfo } from "../api/vpninfo";

export const useVpnStore = defineStore('vpn', () => {
    const vpnInfo = reactive({})
    const getVpnInfo = async () => {
        const res = await getInfo()
        vpnInfo.value = res
    }
    return { vpnInfo, getVpnInfo }
})

