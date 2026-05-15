<template>
  <div>
    <el-carousel :interval="4000" type="card" height="300px">
      <el-carousel-item v-for="item in imgs" :key="item">
        <el-image style="height: 300px" :src="item" fit="contain"></el-image>
      </el-carousel-item>
    </el-carousel>

    <div style="text-align: center; margin-top: 20px;">
      <el-link href="/android" type="primary">Android 文件下载</el-link>
      <el-link href="/ios" type="success" style="margin-left: 20px;">iOS 文件下载</el-link>
    </div>
  </div>
</template>

<script>
import img1 from "../assets/indexpics/blue.jpg";
import img2 from "../assets/indexpics/grass.jpg";
import img3 from "../assets/indexpics/pink.jpg";
import img4 from "../assets/indexpics/purple.jpg";
import img5 from "../assets/indexpics/ranbaw.jpg";
export default {
  name: "Index",
  data() {
    return {
      imgs: [img1, img2, img3, img4, img5],
    };
  },
};
</script>
