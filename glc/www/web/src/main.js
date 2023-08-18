import { createApp } from 'vue';
import ElementPlus from 'element-plus';

import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import ContextMenu from '@imengyu/vue3-context-menu';
import gxui, { pinia } from '~/pkgs';

import App from '~/App.vue';
import { router } from '~/router';

import 'virtual:svg-icons-register';
import SvgIcon from '~/components/svg-icon/SvgIcon.vue';
import GxButton from '~/components/form/GxButton.vue';
import GxUploadButton from '~/components/form/GxUploadButton.vue';

import '~/assets/style/s1-reset.css';
import '~/assets/style/s2-common.css';
import 'element-plus/dist/index.css';
import '~/assets/style/s3-custom.scss';
import '~/assets/style/s4-gxui.scss';

const app = createApp(App);

// 全局注册el-icon
for (const [name, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(name, component);
}

// 全局注册SvgIcon
app.component('SvgIcon', SvgIcon).component('GxButton', GxButton).component('GxUploadButton', GxUploadButton);
app.use(pinia).use(router).use(ElementPlus);
app.use(gxui);
app.use(ContextMenu);

// 自定义指令v-has，控制组件没权限时不显示
app.use({
  install(Vue) {
    Vue.directive('has', {
      mounted(el, binding) {
        !(localStorage.getItem('gotoeasy__permission') || '').split(',').includes(binding.value) && el.parentNode.removeChild(el);
      },
    });
  },
});

app.mount('#app');
