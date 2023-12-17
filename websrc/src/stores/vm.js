import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { getInfo } from "../api/vm";

export const useVpnStore = defineStore('vpn', () => {
    const vmIpInfo = reactive({})
    const ovpnStatusResult = ref("'Status': get status<p/>'Restart': restart openvpn");
    const getVmIpInfo = async () => {
        const res = await getInfo()
        vmIpInfo.value = res
    }
    return { vmIpInfo, ovpnStatusResult, getVmIpInfo }
})

