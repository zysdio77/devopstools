<template>
  <div id="app">
    <!--img src="./assets/logo.png"-->
    <el-header>
      <el-menu
        :default-active="activeIndex"
        class="el-menu-demo"
        mode="horizontal"
        @select="handleSelect"
      >
        <el-menu-item index="/index">主页</el-menu-item>
        <el-submenu index="/download">
          <template slot="title">下载</template>
          <el-menu-item index="/downloadios">下载IOS</el-menu-item>
          <el-menu-item index="/downloadandroid">下载Andriod</el-menu-item>
        </el-submenu>
        <el-submenu index="/upload">
          <template slot="title">包管理</template>
          <el-menu-item index="/uploadios">上传IOS</el-menu-item>
          <el-menu-item index="/uploadandroid">上传Andriod</el-menu-item>
        </el-submenu>
        <el-menu-item index="/collect">收藏</el-menu-item>
      </el-menu>
    </el-header>
    <router-view />
  </div>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      activeIndex: "1",
    };
  },
  methods: {
    handleSelect(key, keyPath) {
      console.log(key, keyPath);
      this.$router.push(key);
    },
  },
};
</script>

<style>
</style>
