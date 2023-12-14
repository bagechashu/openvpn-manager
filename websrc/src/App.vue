<!-- @format -->

<template>
  <header>
    <img src="./assets/image/ovpnmgr.png" style="width: 40px; height: 40px; margin: 12px 10px 0 10px" />
    <h2 style="margin-top: 12px">OpenVPN Manager</h2>
  </header>

  <n-message-provider>
    <message-container />
  </n-message-provider>

  <n-dialog-provider>
    <dialog-container />
  </n-dialog-provider>

  <div class="wrapper">
    <n-tabs type="segment">
      <n-tab-pane name="vpn-info" tab="VPN Info">
        <VpnInfo />
      </n-tab-pane>
      <n-tab-pane name="user-manager" tab="User Manager">
        <UserManager />
      </n-tab-pane>
    </n-tabs>
  </div>

  <footer>&copy; Based on <a href="https://github.com/Nyr/openvpn-install">OpenVPN installation script</a> maintained by Nyr 2023</footer>
</template>

<script setup>
import VpnInfo from "./components/VpnInfo.vue";
import UserManager from "./components/UserManager.vue";
import { NMessageProvider, NDialogProvider } from "naive-ui";
import MessageContainer from "./components/global/MessageContainer.vue";
import DialogContainer from "./components/global/DialogContainer.vue";
</script>

<style scoped>
header {
  display: flex;
  line-height: 1.5;
  margin-bottom: 1rem;
}

@media (min-width: 1024px) {
  header {
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
