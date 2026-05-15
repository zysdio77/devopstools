<template>
  <div>
    <el-main>
      <h1>上传Android</h1>

      <el-upload
        class="upload-demo"
        ref="upload"
        action="/fir/upload"
        name="upload-file"
        :data="uploadData"
        :before-upload="beforeupload"
        :on-preview="handlePreview"
        :on-remove="handleRemove"
        :file-list="fileList"
        :auto-upload="false"
        :on-success="success"
        :on-error="uploadError"
      >
        <el-button slot="trigger" size="small" type="primary">选取文件</el-button>

        <el-button
          style="margin-left: 10px"
          size="small"
          type="success"
          @click="submitUpload"
          :loading="uploading"
        >上传到服务器</el-button>

        <el-input
          name="note"
          type="textarea"
          autosize
          placeholder="请输入备注信息"
          v-model="textarea1"
        >
        </el-input>
        <div slot="tip" class="el-upload__tip">仅支持 .apk 文件，大小不超过500MB</div>
      </el-upload>
    </el-main>
  </div>
</template>

<script>
export default {
  name: "UploadAndroid",
  data() {
    return {
      fileList: [],
      textarea1: '',
      uploading: false,
      uploadData: { note: '', system_type: '', name: '' }
    };
  },
  methods: {
    submitUpload() {
      this.$refs.upload.submit();
    },
    handleRemove(file, fileList) {
      console.log(file, fileList);
    },
    handlePreview(file) {
      console.log(file);
    },
    beforeupload(file) {
      this.uploadData.note = this.textarea1;
      this.uploadData.system_type = "android";
      this.uploadData.name = file.name;
    },
    success(response, file, fileList) {
      if (response.success) {
        this.$message.success('上传成功');
      } else {
        this.$message.error(response.message || '上传失败');
      }
    },
    uploadError(err, file, fileList) {
      this.$message.error('上传失败: ' + err);
    },
  },
};
</script>
