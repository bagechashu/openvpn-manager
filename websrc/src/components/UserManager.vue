<!-- @format -->

<template>
  <div class="controlbar-wrapper">
    <n-h3 prefix="bar" class="controlbar-left">User List</n-h3>
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
import { NDataTable, NInput, NButton, NModal } from "naive-ui";
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

// TODO: 设置隐藏 delete 按钮
const columns = ref([
  {
    title: "Username",
    key: "username",
    minWidth: 100,
  },
  {
    title: "Manage",
    key: "manage",
    align: "center",
    titleAlign: "center",
    minWidth: 200,
    render(row) {
      return h("div", [
        h("span", [
          h("img", {
            src: "/image/CloudDownload.svg",
          }),
          h(
            "a",
            {
              href: `/api/user/cert/${row.username}.ovpn`,
            },
            { default: () => "OvpnFile" }
          ),
        ]),
        h(
          NButton,
          {
            // strong: true,
            // tertiary: true,
            size: "small",
            type: "error",
            style: {
              marginLeft: "1.5rem",
            },
            onClick: () => {
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
            },
          },
          { default: () => "delete" }
        ),
      ]);
    },
  },
]);

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
