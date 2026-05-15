<template>
  <div>
    <el-radio-group v-model="radio1" @input="handleradiobutton()">
      <el-radio-button label="ios"></el-radio-button>
      <el-radio-button label="android"></el-radio-button>
    </el-radio-group>

    <el-table :data="tableData" style="width: 100%" v-loading="loading">
      <el-table-column label="包名">
        <template slot-scope="scope">
          <el-popover trigger="hover" placement="top">
            <p>系统: {{ scope.row.system_type }}</p>
            <p>上传时间: {{ scope.row.create_time }}</p>
            <p>详情: {{ scope.row.note }}</p>
            <div slot="reference" class="name-wrapper">
              <el-tag size="medium">{{ scope.row.name }}</el-tag>
            </div>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button
            type="success"
            @click="handleDownLoad(scope.index, scope.row)"
          >下载</el-button>
          <el-popconfirm
            confirm-button-text="好的"
            cancel-button-text="不用了"
            icon="el-icon-info"
            icon-color="red"
            title="确认删除？"
            @confirm="handleDelete(scope.$index, scope.row)"
          >
            <el-button type="danger" slot="reference">删除</el-button>
          </el-popconfirm>
          <el-popconfirm
            confirm-button-text="是"
            cancel-button-text="否"
            icon="el-icon-info"
            icon-color="blue"
            title="确认取消收藏？"
            @confirm="handleUpdate(scope.$index, scope.row)"
          >
            <el-button type="info" slot="reference">取消收藏</el-button>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="currentPage"
      :page-sizes="[10, 20, 30, 40]"
      :page-size="pageSize"
      :total="totalItems"
    >
    </el-pagination>
  </div>
</template>

<script>
import axios from "axios";
export default {
  data() {
    return {
      tableData: [],
      loading: false,

      currentPage: 1,
      pageSize: 10,
      totalItems: 0,

      radio1: "ios",
      pre: "collect",
    };
  },
  methods: {
    handleDownLoad(index, row) {
      if (this.radio1 == "ios") {
        var plistname = row.name.split(".ipa").join(".plist");
        var addrpath = "/plist/" + plistname;
        var aa = "itms-services:///?action=download-manifest&url=" + window.location.origin + addrpath;
        window.location.href = aa;
      } else {
        window.location.href = "/android/" + row.name;
      }
    },
    handleDelete(index, row) {
      axios.post("/fir/delete", {
        name: row.name,
        system_type: this.radio1
      }).then((res) => {
        if (res.data.success) {
          this.$message.success('删除成功');
          this.loadData();
        } else {
          this.$message.error(res.data.message || '删除失败');
        }
      }).catch(() => {
        this.$message.error('网络错误');
      });
    },
    handleUpdate(index, row) {
      axios.post("/fir/update", {
        name: row.name,
        system_type: this.radio1
      }).then((res) => {
        if (res.data.success) {
          this.$message.success('已取消收藏');
          this.loadData();
        } else {
          this.$message.error(res.data.message || '操作失败');
        }
      }).catch(() => {
        this.$message.error('网络错误');
      });
    },
    handleradiobutton() {
      this.currentPage = 1;
      this.loadData();
    },
    handleSizeChange(val) {
      this.pageSize = val;
      this.currentPage = 1;
      this.loadData();
    },
    handleCurrentChange(val) {
      this.currentPage = val;
      this.loadData();
    },
    loadData() {
      this.loading = true;
      axios.get("/fir/page", {
        params: { name: this.pre + this.radio1, page: this.currentPage }
      }).then((res) => {
        this.tableData = res.data.result || [];
        this.totalItems = res.data.total || 0;
      }).catch(() => {
        this.$message.error('加载数据失败');
      }).finally(() => {
        this.loading = false;
      });
    },
  },
  created() {
    this.loadData();
  },
};
</script>
