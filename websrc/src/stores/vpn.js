import { defineStore } from 'pinia'
import { reactive } from 'vue'
import { getInfo } from "../api/vpninfo";

export const useVpnStore = defineStore('vpn', () => {
    const vpnInfo = reactive({})
    async function getVpnInfo() {
        getInfo()
            .then((res) => {
                vpnInfo.value = res;
            })
            .catch((error) => {
                console.error("Error fetching VPN info:", error);
            });
    }
    return { vpnInfo, getVpnInfo }
})

