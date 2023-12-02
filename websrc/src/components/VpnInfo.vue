<!-- @format -->
<template>
  <n-list bordered>
    <n-list-item>
      <n-thing v-for="(value, name, key) in vpnData" :key="key">
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

<script>
import { getInfo } from "../api/vpninfo";

export default {
  name: "VpnInfo",
  data() {
    return {
      vpnData: {},
    };
  },
  async created() {
    this.vpnData = await this.getVpnInfo();
    // console.log(this.vpnData);
  },

  methods: {
    async getVpnInfo() {
      try {
        return await getInfo();
      } catch (error) {
        console.error("Error fetching VPN info:", error);
        return {};
      }
    },
  },
};
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
  white-space:nowrap;
  word-break:keep-all;
  text-overflow:ellipsis;
  overflow:hidden;
}

</style>
