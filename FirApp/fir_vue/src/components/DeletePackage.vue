<template>
  <div>
    <el-radio-group v-model="radio1">
      <el-radio-button label="ios"></el-radio-button>
      <el-radio-button label="android"></el-radio-button>
    </el-radio-group>

    <el-button type="primary" icon="el-icon-search" @click="gettableData">搜索</el-button>

    <el-table
      ref="multipleTable"
      :data="tableData"
      tooltip-effect="dark"
      style="width: 100%"
      v-loading="loading"
      @selection-change="handleSelectionChange">
      <el-table-column
        type="selection"
        width="55">
      </el-table-column>
      <el-table-column
        prop="name"
        label="包名"
        width="120">
      </el-table-column>
      <el-table-column
        prop="system_type"
        label="系统"
        width="120">
      </el-table-column>
      <el-table-column
        prop="note"
        label="备注"
        show-overflow-tooltip>
      </el-table-column>
      <el-table-column
        fixed="right"
        label="操作"
        width="120">
        <template slot-scope="scope">
          <el-button
            @click.native.prevent="deleteRow(scope.$index, scope.row)"
            type="text"
            size="small">
            移除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <div style="margin-top: 10px">
      <el-button type="danger" @click="batchDelete" :disabled="multipleSelection.length === 0">
        批量删除 ({{ multipleSelection.length }})
      </el-button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
export default {
  data() {
    return {
      radio1: 'ios',
      tableData: [],
      loading: false,
      multipleSelection: []
    }
  },
  methods: {
    deleteRow(index, row) {
      this.$confirm('确认删除该包？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        axios.post("/fir/delete", {
          name: row.name,
          system_type: row.system_type
        }).then((res) => {
          if (res.data.success) {
            this.$message.success('删除成功');
            this.tableData.splice(index, 1);
          } else {
            this.$message.error(res.data.message || '删除失败');
          }
        }).catch(() => {
          this.$message.error('网络错误');
        });
      }).catch(() => {});
    },
    toggleSelection(rows) {
      if (rows) {
        rows.forEach(row => {
          this.$refs.multipleTable.toggleRowSelection(row);
        });
      } else {
        this.$refs.multipleTable.clearSelection();
      }
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    batchDelete() {
      if (this.multipleSelection.length === 0) return;
      this.$confirm(`确认删除选中的 ${this.multipleSelection.length} 个包？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let promises = this.multipleSelection.map(row => {
          return axios.post("/fir/delete", {
            name: row.name,
            system_type: row.system_type
          });
        });
        Promise.all(promises).then(() => {
          this.$message.success('批量删除完成');
          this.gettableData();
        }).catch(() => {
          this.$message.error('部分删除失败');
          this.gettableData();
        });
      }).catch(() => {});
    },
    gettableData() {
      this.loading = true;
      axios.get("/fir/info", {
        params: { name: this.radio1 }
      }).then((res) => {
        this.tableData = Array.isArray(res.data) ? res.data : JSON.parse(res.data);
      }).catch(() => {
        this.$message.error('加载数据失败');
      }).finally(() => {
        this.loading = false;
      });
    },
  }
}
</script>
