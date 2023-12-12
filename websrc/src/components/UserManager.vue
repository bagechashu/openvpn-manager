<!-- @format -->

<template>
  <div class="controlbar-wrapper">
    <n-h3 prefix="bar" class="controlbar-left">User List</n-h3>
    <div class="controlbar-right">
      <n-button @click="showCreateUserModal = true" size="small" color="#44aaee"> Create </n-button>
      <n-button size="small" type="warning"> Edit </n-button>
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
    <n-input v-model:value="username" type="text" placeholder="ovpn-20231212" />
  </n-modal>
  <n-data-table :columns="columns" :bordered="false" :data="data" :pagination="pagination" />
  <n-back-top :right="100" />
</template>

<script setup>
import { ref, h, defineComponent } from "vue";
import { NButton, NModal } from "naive-ui";

const data = [{ username: "Wonderwall" }, { username: "Don't Look Back in Anger" }, { username: "Champagne Supernova" }];
const pagination = {
  pageSize: 10,
};
const createColumns = ({ deleteUser }) => {
  return [
    {
      title: "Username",
      key: "username",
    },
    {
      title: "Download",
      key: "download",
      render(row) {
        return h(
          "a",
          {
            href: `/api/user/cert/${row.username}.ovpn`,
          },
          { default: () => "Download Ovpn" }
        );
      },
    },
    {
      title: "Delete",
      key: "delete",
      render(row) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: "small",
            onClick: () => deleteUser(row),
          },
          { default: () => "delete" }
        );
      },
    },
  ];
};
const columns = createColumns({
  deleteUser(row) {
    window.$message.info(`deleteUser ${row.title}`);
  },
});

const showCreateUserModal = ref(false);
let username = "";
function cancelCallback() {
  window.$message.warning("Cancel");
}
function submitCallback() {
  window.$message.success("Submit");
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
.table-control.n-button {
  margin-right: 0.5rem;
}
.n-button {
  width: 4rem;
}
</style>
