<!-- @format -->
<template>
  <n-list bordered>
    <n-list-item>
      <n-thing v-for="(value, name, key) in vmIpInfo.value" :key="key">
        <n-grid :cols="16" item-responsive>
          <n-grid-item class="name" span="0 500:3">
            {{ name.toUpperCase() }}
          </n-grid-item>
          <n-grid-item class="value" span="16 500:13">
            {{ value }}
          </n-grid-item>
        </n-grid>
      </n-thing>
    </n-list-item>
  </n-list>
  <n-card title="Openvpn Service Status">
    <template #header-extra>
      <n-button @click="clickOvpnStatus" size="small" type="info"> Check </n-button>
    </template>
    <div v-html="ovpnStatusResult"></div>
  </n-card>
</template>

<script setup>
import { NGrid, NGridItem, NList, NListItem, NThing, NButton, NCard } from "naive-ui";
import { storeToRefs } from "pinia";
import { useVpnStore } from "../stores/vm";
import { getInfo, ovpnStatus } from "../api/vm";

const vpnStore = useVpnStore();
const { vmIpInfo, ovpnStatusResult } = storeToRefs(vpnStore);

if (Object.keys(vmIpInfo.value).length === 0) {
  getInfo()
    .then((res) => {
      vpnStore.vmIpInfo.value = res;
      // console.log(vmIpInfo.value);
    })
    .catch((error) => {
      console.error("Error fetching VPN info:", error);
    });
}

// onMounted(() => {
//   if (Object.keys(vmIpInfo.value).length !== 0) return;
//   vpnStore.getVmIpInfo();
// });

function clickOvpnStatus() {
  ovpnStatus()
    .then((res) => {
      // console.log(res.response.msg);
      vpnStore.ovpnStatusResult = `<pre class="scrollable-pre">${res.response.msg}</pre>`;
    })
    .catch((error) => {
      console.log(error);
      window.$message.error(`get openvpn status error: ${error}`, { duration: 2000 });
    });
}
</script>

<style scope>
.name {
  display: flex;
  align-items: center;
  height: 2.5rem;
  color: #54bedd;
  font-weight: bold;
}
.value {
  display: flex;
  align-items: center;
  height: 2.5rem;
  white-space: nowrap;
  word-break: keep-all;
  text-overflow: ellipsis;
  overflow: hidden;
}
.n-card {
  margin-top: 1rem;
}
.scrollable-pre {
  white-space: pre;
  overflow-x: auto;
  font-size: 0.6rem;
}
</style>
