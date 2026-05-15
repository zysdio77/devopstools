# FirApp - 前端

移动应用分发平台的 Web 管理界面。

## 技术栈

Vue 2.x + Element UI + Vue Router + Axios

## 目录结构

```
src/components/
├── Index.vue              # 首页
├── UploadIOS.vue          # 上传 iOS 包
├── UploadAndroid.vue      # 上传 Android 包
├── DownloadIOS.vue        # iOS 下载列表
├── DownloadAndroid.vue    # Android 下载列表
├── Collect.vue            # 收藏列表
└── DeletePackage.vue      # 批量删除
```

## 本地开发

```bash
npm install
npm run dev        # 开发模式 localhost:8080
npm run build      # 构建到 dist/，拷贝到 ../fir_go/template/dist/
```

## 与后端通信

所有请求使用相对路径 `/fir/*`，部署时前后端同域无需额外配置。
