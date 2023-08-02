<script setup>
import {ref} from "vue";
import {ElNotification} from "element-plus";
import {DbInfo} from "../../wailsjs/go/main/App.js";

let dialogVisible = ref(false)
let props = defineProps(['data'])
let info = ref()

function getDBInfo() {
  dialogVisible.value = true
  DbInfo(props.data.identity).then(res => {
    if (res.code != 200) {
      ElNotification({
        title:res.msg,
        type: "error",
      })
      return
    }
    info.value = res.data
  });
}

</script>

<template>
<main>
  <el-button type="primary" @click="getDBInfo" link>详情</el-button>
  <el-dialog
      v-model="dialogVisible"
      title="数据库详情"
      width="60%"
  >
    <el-table :data="info" border stripe style="width: 100%">
      <el-table-column type="index" width="65"/>
      <el-table-column prop="key" label="Key" />
      <el-table-column prop="value" label="Value" />
    </el-table>
  </el-dialog>
</main>
</template>

<style scoped>

</style>