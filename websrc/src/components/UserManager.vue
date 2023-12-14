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
    style="position: fixed; top: 6rem; left: 50%; transform: translateX(-50%)"
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
  <n-back-top :right="100" />
</template>

<script setup>
import { ref, h } from "vue";
import { storeToRefs } from "pinia";
import { NDataTable, NInput, NButton, NModal, NSwitch, NPopconfirm } from "naive-ui";
import { useUsersStore } from "../stores/users";
import { getUsers, createUser, deleteUser } from "../api/user";

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
    minWidth: 100,
    render(row) {
      return h("div", [
        h(
          "a",
          {
            href: `/api/user/cert/${row.username}.ovpn`,
          },
          { default: () => "OvpnFile" }
        ),
        h("img", { src: "/image/CloudDownload.svg" }),
        h(
          NPopconfirm,
          {
            onPositiveClick() {
              delUser(row);
            },
            onNegativeClick() {
              window.$message.warning(`cancel`);
            },
            negativeText: "cancel",
            positiveText: "confirm",
            showIcon: true,
          },
          {
            trigger: () =>
              h(
                NButton,
                {
                  strong: true,
                  size: "small",
                  type: "error",
                  disabled: !isEdit.value,
                },
                { default: () => "delete" }
              ),
            default: () => h("span", {}, { default: () => `Are you sure to delete [${row.username}]?` }),
          }
        ),
      ]);
    },
  },
]);

function delUser(row) {
  deleteUser(row.username)
    .then((res) => {
      console.log(res.response.msg);
      // location.reload();
      window.$message.info(`[${row.username}] is deleted`);
      // 删除用户后, 刷新 users
      usersStore.delUserLocal(row.username);
      // console.log(usersStore.users.value);
    })
    .catch((error) => {
      window.$message.error("Error while deleting users:", error);
    });
}

const showCreateUserModal = ref(false);
const username = ref("");
function cancelCallback() {
  window.$message.warning("Cancel");
}
function submitCallback() {
  const data = JSON.stringify({ username: username.value });
  createUser(data)
    .then((res) => {
      console.log(res.response.msg);
      // 创建用户后, 刷新 users
      usersStore.addUserLocal({ username: username.value });
      // location.reload();
      window.$message.success(`[${username.value}] is created`);
    })
    .catch((error) => {
      window.$message.error("Error when creating users:", error);
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
.n-button {
  width: 4rem;
}
</style>
