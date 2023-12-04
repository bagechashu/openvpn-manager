<!-- @format -->
<template>
  <n-list bordered>
    <n-list-item>
      <n-thing v-for="(value, name, key) in vpnInfo.value" :key="key">
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
</template>

<script setup>
import { storeToRefs } from "pinia";
import { useVpnStore } from "../stores/vpn";
import { getInfo } from "../api/vpninfo";

const vpnStore = useVpnStore();
const { vpnInfo } = storeToRefs(vpnStore);

if (Object.keys(vpnInfo.value).length === 0) {
  getInfo()
    .then((res) => {
      vpnStore.vpnInfo.value = res;
      console.log(vpnInfo.value);
    })
    .catch((error) => {
      console.error("Error fetching VPN info:", error);
    });
}

// onMounted(() => {
//   if (Object.keys(vpnInfo.value).length !== 0) return;
//   vpnStore.getVpnInfo();
// });
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
</style>
