import nProgress from 'nprogress';
import 'nprogress/nprogress.css';

// 路由跳转时的进度条显示配置
nProgress.configure({
  // easing: 'ease', // 动画方式
  speed: 100, // 递增进度条的速度
  showSpinner: false, // 是否显示加载ico
  // trickleSpeed: 200, // 自动递增间隔
  // minimum: 0.3, // 初始化时的最小百分比
});

export default nProgress;
