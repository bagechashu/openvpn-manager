import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { getInfo } from "../api/vm";

export const useVpnStore = defineStore('vpn', () => {
    const vpnInfo = reactive({})
    const ovpnStatusResult = ref("click 'check' to get!");
    const getVpnInfo = async () => {
        const res = await getInfo()
        vpnInfo.value = res
    }
    return { vpnInfo, ovpnStatusResult, getVpnInfo }
})

