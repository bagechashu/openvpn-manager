<!-- @format -->

<template>
  <div class="controlbar-wrapper">
    <div class="controlbar-left"><span style="margin-right: 5px">Edit</span><n-switch size="small" v-model:value="isEdit" /></div>
    <div class="controlbar-right">
      <n-button @click="showCreateUserModal = true" size="small" color="#44aaee"> Create </n-button>
    </div>
  </div>
  <n-modal
    v-model:show="showCreateUserModal"
    preset="dialog"
    title="create user"
    positive-text="submit"
    negative-text="cancel"
    style="position: fixed; top: 11rem; left: 50%; transform: translateX(-50%)"
    @positive-click="submitCallback"
    @negative-click="cancelCallback">
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
  <n-back-top :bottom="5" :right="5" :visibility-height="300"/>
</template>

<script setup>
import { ref, h } from "vue";
import { storeToRefs } from "pinia";
import { NDataTable, NInput, NButton, NModal, NSwitch, NIcon, NBackTop } from "naive-ui";
import { useUsersStore } from "../stores/users";
import { getUsers, createUser, deleteUser } from "../api/user";
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

// const users = [{ username: "Wonderwall" }, { username: "Don't Look Back in Anger" }, { username: "Champagne Supernova" }];
const pagination = {
  pageSize: 10,
};

const isEdit = ref(false);
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
            onClick: () => downloadCert(row.username),
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
            onClick: () => delAfterConfirm(row),
          },
          { default: () => "delete" }
        ),
      ]);
    },
  },
]);
function downloadCert(username) {
  axios({
    method: "get",
    url: `/api/user/cert/${username}.ovpn`,
    responseType: "blob",
  })
    .then((res) => {
      let data = res.data; // 这里后端对文件流做了一层封装，将data指向res.data即可
      if (!data) {
        return;
      }
      let url = window.URL.createObjectURL(new Blob([data]));
      let a = document.createElement("a");
      a.style.display = "none";
      a.href = url;
      a.setAttribute("download", `${username}.ovpn`);
      document.body.appendChild(a);
      a.click(); //执行下载
      window.URL.revokeObjectURL(a.href); //释放url
      document.body.removeChild(a); //释放标签
    })
    .catch((error) => {
      // console.log(error);
      window.$message.error(`${error.response.status} ${error.response.statusText}`, { duration: 2000 });
    });
}
function delAfterConfirm(row) {
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
          window.$message.success(`[ ${row.username} ] is deleted`, { duration: 2000 });
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

const showCreateUserModal = ref(false);
const username = ref("");
function cancelCallback() {
  window.$message.warning("Canceled!");
}
function submitCallback() {
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
