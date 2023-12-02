<!-- @format -->

<template>
  <div></div>
</template>

<script>
// https://blog.csdn.net/weixin_43856797/article/details/124377784
import { defineComponent } from "vue";
import { useMessage } from "naive-ui";

export default defineComponent({
  setup() {
    window.$message = useMessage();
  },
});
</script>
