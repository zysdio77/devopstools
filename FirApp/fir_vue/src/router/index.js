import Vue from 'vue'
import Router from 'vue-router'
//import HelloWorld from '@/components/HelloWorld'
import Index from '@/components/Index';
import DownloadIos from '@/components/DownloadIOS';
import DownloadAndroid from '@/components/DownloadAndroid';
import UploadAndroid from '@/components/UploadAndroid';
import UploadIOS from '@/components/UploadIOS';
import Collect from '@/components/Collect';


Vue.use(Router)

export default new Router({
  routes: [
    {
      //路由重定向
      path: '/',
      redirect: '/index'
      //name: 'HelloWorld',
      //component: HelloWorld
    },
    {
      path: '/index',
      name: 'Index',
      component: Index
    },
    {
      path: '/downloadios',
      name: 'DownloadIos',
      component: DownloadIos
    },
    {
      path: '/downloadandroid',
      name: 'DownloadAndroid',
      component: DownloadAndroid
    },
    {
      path: '/uploadandroid',
      name: 'UploadAndroid',
      component: UploadAndroid
    },
    {
      path: '/uploadios',
      name: 'UploadIOS',
      component: UploadIOS
    },
    {
      path: '/collect',
      name: 'Collect',
      component: Collect
    },

  ]
})
