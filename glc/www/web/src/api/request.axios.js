import axios from 'axios';
import { config } from '~/config';
import { $msg, $emitter, useTokenStore } from '~/pkgs';
import { mockUrls } from '~/config/mock-urls.config';

let tokenStore = null;
const { VITE_AXIOS_BASE_URL, VITE_DEV_DOMAIN_URL, VITE_API_DOMAIN_PATH, VITE_DEV_MOCK } = import.meta.env;

const msgSet = new Set();

// 创建请求实例
const axiosInstance = axios.create({
  baseURL: VITE_AXIOS_BASE_URL,
  timeout: config.TIMEOUT, // 指定请求超时的毫秒数
});

function getRequestUrl(url) {
  if (/^http/i.test(url)) {
    return url; // 绝对地址时直接返回
  }

  let { origin } = window.location;
  if (/^http:\/\/127.0.0.1/i.test(origin) || (/^http:\/\/localhost/i.test(origin) && VITE_DEV_DOMAIN_URL)) {
    origin = VITE_DEV_DOMAIN_URL; // 本地开发模式时使用配置
  }

  if (VITE_API_DOMAIN_PATH && VITE_API_DOMAIN_PATH !== './' && VITE_API_DOMAIN_PATH !== '""' && VITE_API_DOMAIN_PATH !== "''") {
    origin += VITE_API_DOMAIN_PATH;
  }

  const rs = origin + url; // 非本地开发模式时使用浏览器地址
  // console.log('[request]', rs);
  return rs;
}

// 前置拦截器（发起请求之前的拦截）
axiosInstance.interceptors.request.use(
  conf => {
    // 修改请求地址（配为mock时保持不变以让其被拦截走mock，否则调整url）
    if (!VITE_DEV_MOCK || !mockUrls.includes(conf.url)) {
      conf.url = getRequestUrl(conf.url);
    }

    // 发送请求之前判断加入令牌消息头
    !tokenStore && (tokenStore = useTokenStore());
    if (tokenStore.needLogin == 'true') {
      conf.headers[config.TOKEN_HEADER] = tokenStore.token;
    } else {
      delete conf.headers[config.TOKEN_HEADER];
    }
    console.log('axiosInstance-conf', conf.url);
    return conf;
  },
  // 错误处理
  error => Promise.reject(error),
);

// 后置拦截器（获取到响应时的拦截）
axiosInstance.interceptors.response.use(
  response => response,
  error => {
    console.log(error);
    const { response, message } = error;
    if (response?.status == '404') {
      $msg.notify(message, 'error');
      return Promise.reject(error);
    }

    if (response?.status == '502') {
      $msg.notify('无效的网关（502）', 'error');
      return Promise.reject(error);
    }

    if (message == 'Network Error') {
      return $msg.notify('网络错误！', 'error');
    }

    if (message == '会话已超时，请重新登录!') {
      !tokenStore && (tokenStore = useTokenStore());
      // 4小时内令牌，且1秒内未重复报过才允许报消息
      if (new Date().getTime() - tokenStore.time < 4 * 60 * 60 * 1000 && !msgSet.has(data.message)) {
        msgSet.add(data.message);
        setTimeout(() => msgSet.delete(data.message), 1000); // 1秒后才允许报重复消息
        $msg.notify(data.message, 'error');
      }
      return;
    }

    const data = response?.data || {};
    if (data.code == '401') {
      !tokenStore && (tokenStore = useTokenStore());
      // 4小时内令牌，且1秒内未重复报过才允许报消息
      if (new Date().getTime() - tokenStore.time < 4 * 60 * 60 * 1000 && !msgSet.has(data.message)) {
        msgSet.add(data.message);
        setTimeout(() => msgSet.delete(data.message), 1000); // 1秒后才允许报重复消息
        $msg.notify(data.message, 'error');
      }
      return;
    }

    if (response?.data) {
      return Promise.reject(error);
    }

    $msg.notify(message, 'error');

    return Promise.reject(error);
  },
);

// 导出常用函数

/**
 * post请求
 *
 * @param {string} url 请求地址
 * @param {object} data 请求数据
 * @param {object} params 请求地址中的参数
 * @param {object} headers 请求头
 */
export const postAction = (url, data = {}, params = {}, headers = { 'Content-Type': 'application/json;charset=UTF-8' }) => {
  !tokenStore && (tokenStore = useTokenStore());
  if (tokenStore.needLogin == 'true') {
    data.token = tokenStore.token;
  }
  return new Promise(resolve => {
    axiosInstance({
      method: 'post',
      url,
      data,
      params,
      headers,
    })
      .then(rs => {
        rs.data ? resolve(rs.data) : resolve(rs); // rs.data对应后端特定接口
      })
      .catch(e => {
        resolve({ error: e }); // 统一处理请求异常，后置拦截器已报错，这里直接忽略
      });
  });
};

/**
 * get请求
 *
 * @param {string} url 请求地址
 * @param {object} params 请求地址中的参数
 * @param {object} headers 请求头
 */
export const getAction = (url, params = {}, headers = { 'Content-Type': 'application/json;charset=UTF-8' }) =>
  new Promise(resolve => {
    axiosInstance({
      method: 'get',
      url,
      params,
      headers,
    })
      .then(rs => {
        rs.data ? resolve(rs.data) : resolve(rs); // rs.data对应后端特定接口
      })
      .catch(e => {
        resolve({ error: e }); // 统一处理请求异常，后置拦截器已报错，这里直接忽略
      });
  });

/**
 * upload请求
 *
 * @param {string} url 请求地址
 * @param {object} data 请求数据
 * @param {object} params 请求地址中的参数
 */
export const uploadAction = (url, data = {}, params = {}) =>
  new Promise(resolve => {
    axiosInstance({
      method: 'post',
      url,
      data,
      params,
      headers: {
        'Content-Type': 'multipart/form-data', // 文件上传
      },
    })
      .then(rs => {
        rs.data ? resolve(rs.data) : resolve(rs); // rs.data对应后端特定接口
      })
      .catch(e => {
        resolve({ error: e }); // 统一处理请求异常，后置拦截器已报错，这里直接忽略
      });
  });

// ------------------------------------------------
export const getUrl = getRequestUrl;
window.$getUrl = getUrl;

// 以下全局自动import，供直接使用
export const $post = postAction;
export const $get = getAction;
export const $upload = uploadAction;

$emitter.on('$axios:post', postAction);
$emitter.on('$axios:get', getAction);
$emitter.on('$axios:upload', uploadAction);
