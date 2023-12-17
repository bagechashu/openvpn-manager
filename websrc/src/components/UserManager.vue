<!-- @format -->

<template>
  <div class="controlbar-wrapper">
    <div class="controlbar-left"><span style="margin-right: 5px">Edit</span><n-switch size="small" v-model:value="isEdit" /></div>
    <div class="controlbar-right">
      <n-button v-show="!isEdit" @click="showUserCreateModal = true" size="small" color="#44aaee"> Create </n-button>
      <n-button v-show="isEdit" @click="clickOvpnRestart" size="small" type="warning"> restart </n-button>
    </div>
  </div>
  <n-modal
    v-model:show="showUserCreateModal"
    preset="dialog"
    title="create user"
    positive-text="submit"
    negative-text="cancel"
    style="position: fixed; top: 11rem; left: 50%; transform: translateX(-50%)"
    @positive-click="clickUserCreate"
    @negative-click="clickCancelUserCreate">
    <div>
      <div>username</div>
      <n-input v-model:value="username" type="text" placeholder="eg: ovpn-20231212" />
    </div>
  </n-modal>
  <n-data-table
    :columns="columns"
    :bordered="false"
    :single-line="false"
    :striped="true"
    :scroll-x="300"
    :data="usersStore.users.value"
    :pagination="pagination"
    :paginate-single-page="false" />
  <n-back-top :bottom="5" :right="5" :visibility-height="300" />
</template>

<script setup>
import { ref, h } from "vue";
import { storeToRefs } from "pinia";
import { NDataTable, NInput, NButton, NModal, NSwitch, NIcon, NBackTop } from "naive-ui";
import { useUsersStore } from "../stores/users";
import { getUsers, createUser, deleteUser, ovpnRestart } from "../api/user";
import { CloudDownload } from "@vicons/tabler";
import axios from "axios";

const usersStore = useUsersStore();
const { users } = storeToRefs(usersStore);

if (Object.keys(users.value).length === 0) {
  getUsers()
    .then((res) => {
      usersStore.users.value = res.response.users;
      // console.log(res.response.users);
      // console.log(users.value);
    })
    .catch((error) => {
      console.error("Error fetching users:", error);
    });
}

const isEdit = ref(false);

const showUserCreateModal = ref(false);
const username = ref("");
function clickCancelUserCreate() {
  window.$message.warning("Canceled!");
}
function clickUserCreate() {
  const data = JSON.stringify({ username: username.value });
  createUser(data)
    .then((res) => {
      console.log(res.response.msg, res.status);
      if (res.status === "success") {
        // 创建用户后, 刷新 users
        usersStore.addUserLocal({ username: username.value });
        // location.reload();
        window.$message.success(`[${username.value}] creat success`, { duration: 2000 });
      }
    })
    .catch((error) => {
      window.$message.error(`Error when creating users: ${error}`, { duration: 2000 });
    });
}

// const users = [{ username: "Wonderwall" }, { username: "Don't Look Back in Anger" }, { username: "Champagne Supernova" }];
const pagination = {
  pageSize: 10,
};

function clickOvpnRestart() {
  window.$dialog.warning({
    title: "Warning",
    content: `Are you sure to restart openvpn service?`,
    positiveText: "Sure",
    negativeText: "Cancel",
    style: "position: fixed; top: 11rem; left: 50%; transform: translateX(-50%)",
    onPositiveClick: () => {
      ovpnRestart()
        .then((res) => {
          console.log(res.response.msg);
          window.$message.success(`openvpn service is restart`, { duration: 2000 });
        })
        .catch((error) => {
          window.$message.error(`openvpn service restart error: ${error}`, { duration: 2000 });
        });
    },
  });
}

const columns = ref([
  {
    title: "User",
    key: "username",
    align: "center",
    titleAlign: "center",
    minWidth: 50,
  },
  {
    title: "Manage",
    key: "manage",
    align: "center",
    titleAlign: "center",
    minWidth: 200,
    render(row) {
      return h("div", [
        h(
          NButton,
          {
            text: true,
            size: "small",
            strong: true,
            type: "primary",
            style: { display: !isEdit.value ? "inline-flex" : "none" },
            onClick: () => clickCertDownload(row.username),
          },
          {
            default: () => "Download Ovpn File",
            icon: () => h(NIcon, { color: "#000", size: 24 }, { default: () => h(CloudDownload) }),
          }
        ),
        h(
          NButton,
          {
            strong: true,
            size: "small",
            type: "error",
            style: { display: isEdit.value ? "inline-flex" : "none" },
            onClick: () => clickUserDelete(row),
          },
          { default: () => "delete" }
        ),
      ]);
    },
  },
]);
function clickCertDownload(username) {
  axios({
    method: "get",
    url: `/api/user/cert/${username}.ovpn`,
    responseType: "blob",
  })
    .then((res) => {
      let data = res.data; // axios.response.data
      if (!data) {
        return;
      }
      // To let android don't fix extname of file, Blob set type to octet.
      let url = window.URL.createObjectURL(new Blob([data], { type: "application/octet-stream" }));
      let a = document.createElement("a");
      a.style.display = "none";
      a.href = url;
      a.setAttribute("download", `${username}.ovpn`);
      document.body.appendChild(a);
      a.click(); //exec download
      window.URL.revokeObjectURL(a.href); //release url
      document.body.removeChild(a); //release tag
    })
    .catch((error) => {
      // console.log(error);
      window.$message.error(`${error.response.status} ${error.response.statusText}`, { duration: 2000 });
    });
}
function clickUserDelete(row) {
  window.$dialog.warning({
    title: "Warning",
    content: `Are you sure to delete user [ ${row.username} ]?`,
    positiveText: "Sure",
    negativeText: "Cancel",
    style: "position: fixed; top: 11rem; left: 50%; transform: translateX(-50%)",
    onPositiveClick: () => {
      deleteUser(row.username)
        .then((res) => {
          console.log(res.response.msg);
          // location.reload();
          window.$message.success(`[ ${row.username} ] is deleted, click [restart] will finally take effect.`, { duration: 2000 });
          // 删除用户后, 刷新 users
          usersStore.delUserLocal(row.username);
          // console.log(usersStore.users.value);
        })
        .catch((error) => {
          window.$message.error(`Error while deleting users: ${error}`, { duration: 2000 });
        });
    },
    // onNegativeClick: () => {
    //   window.$message.error("Canceled!");
    // },
  });
}
</script>

<style scoped>
.controlbar-wrapper {
  position: relative;
  justify-content: space-between;
  align-items: center;
  height: 2.5rem;
}
.controlbar-left {
  display: inline-block;
  margin-left: 0.5rem;
}
.controlbar-right {
  display: inline-block;
  position: absolute;
  right: 0px;
  margin-right: 0.5rem;
}
</style>
