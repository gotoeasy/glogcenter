import { defineConfig, loadEnv } from 'vite';
import { resolve } from 'path';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers';
import OptimizationPersist from 'vite-plugin-optimize-persist';
import PkgConfig from 'vite-plugin-package-config';
// import importToCDN from 'vite-plugin-cdn-import';
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons';
import { viteMockServe } from 'vite-plugin-mock';

export default ({ mode }) => {
  const { VITE_PORT, VITE_WEB_PATH, VITE_PRD_DROP_DEBUGGER, VITE_DEV_MOCK } = loadEnv(mode, process.cwd());

  return defineConfig({
    base: VITE_WEB_PATH,
    plugins: [
      vue(),
      vueJsx(),
      PkgConfig(),
      OptimizationPersist(),
      AutoImport({
        imports: ['vue', 'vue-router', { '~/utils': ['$post', '$get'] }, { '~/pkgs': ['$msg', '$emitter'] }],
        eslintrc: {
          enabled: false, // 若没'./.eslintrc-auto-import.json'文件，先开启，生成后在关闭
        },
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
      createSvgIconsPlugin({
        iconDirs: [`${resolve(__dirname, 'src/assets/icons')}`], // SVG图标文件夹
        symbolId: 'gx_icon-[name]', // 指定symbolId格式
      }),
      viteMockServe({
        mockPath: './src/mock', // mock文件位置
        localEnabled: VITE_DEV_MOCK, // 是否应用于本地
        watchFiles: true, // 监听mock文件变化
        logger: true, // 开启日志
      }),
    ],
    resolve: {
      alias: {
        '~/': `${resolve(__dirname, 'src')}/`,
      },
    },
    server: {
      https: false, // 是否开启 https
      port: VITE_PORT, // 端口号
      host: '0.0.0.0', // 监听所有地址
      open: true, // 服务启动时是否自动打开浏览器
      cors: true, // 允许跨域
      proxy: {
        // '/boot': {
        //   target: 'http://127.0.0.1:8080/boot',
        //   changeOrigin: true,
        //   rewrite: path => path.replace(/^\/boot/, ''),
        // },
      },
    },
    esbuild: {
      // 打包时要去除的语句
      pure: VITE_PRD_DROP_DEBUGGER ? ['console.log', 'console.debug', 'console.info', 'debugger'] : ['console.log', 'console.debug', 'console.info'],
    },

    build: {
      target: 'es2015', // 设置最终构建的浏览器兼容目标(es2015/modules)
      sourcemap: false, // 构建后是否生成 source map 文件
      chunkSizeWarningLimit: 2000, //  chunk 大小警告的限制（以 kbs 为单位）
      reportCompressedSize: false, // 启用/禁用 gzip 压缩大小报告
      rollupOptions: {
        output: {
          // 打包时按类型生成整洁的文件
          chunkFileNames: 'assets/glc-chunk-[hash].js',
          entryFileNames: 'assets/glc-entry-[hash].js',
          assetFileNames: 'assets/glc-asset-[hash].[ext]',
        },
      },
    },
  });
};
