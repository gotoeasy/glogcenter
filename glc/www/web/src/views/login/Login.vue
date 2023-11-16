<template>
  <div class="login-container flx-center">
    <div class="login-box">
      <div class="login-form">
        <div class="login-logo">
          <!-- <img class="login-icon" src="/image/glc.png" alt="" /> -->
          <img src="/image/title.png" alt="" style="width: 420px;" />
          <!-- <h2 class="logo-text">ADMIN LOGIN</h2> -->
        </div>
        <el-form ref="form" :model="formData" :rules="rules" size="large">
          <el-form-item prop="username">
            <el-input v-model="formData.username" placeholder="请输入用户名" maxlength="100">
              <template #prefix>
                <el-icon class="el-input__icon">
                  <user />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="formData.password" type="password" placeholder="请输入密码" autocomplete="new-password"
              maxlength="100">
              <template #prefix>
                <el-icon class="el-input__icon">
                  <lock />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
        <div class="login-btn">
          <el-button round size="large" type="primary" :loading="loading" @click="login">
            <SvgIcon name="user" />
            <span>登 录</span>
          </el-button>
          <el-button size="large" round @click="resetForm">
            <SvgIcon name="refresh-left" />
            <span>重 置</span>
          </el-button>
        </div>

      </div>
    </div>
    <div class="login-footer">Copyright © 2022-present gotoeasy.top</div>
  </div>
</template>

<script setup>
import { userLogin, enableLogin } from '~/api';
import { useTokenStore } from '~/pkgs';

const route = useRoute();
const router = useRouter();
const tokenStore = useTokenStore();

const login = async () => {
  if (loading.value) return;

  await form.value.validate(async (valid) => {
    if (!valid) {
      return;
    }

    const rs = await userLogin({ ...formData.value });
    if (!rs) {
      return;
    }

    // 跳转到目标页面地址（含参数），未指定目标地址时跳转首页
    const redirect = decodeURIComponent(route.query.redirect || '/');
    console.log("redirect", redirect)
    let path = redirect;
    const query = {};
    const idx = redirect.indexOf('?');
    if (idx > 0) {
      path = redirect.substring(0, idx);
      const nameValueArr = redirect.substring(idx + 1).split("&"); // 多参数
      for (let i = 0; i < nameValueArr.length; i++) {
        const pos = nameValueArr[i].indexOf("=");
        if (pos < 0) continue; // 如果没有找到就跳过
        const argName = nameValueArr[i].substring(0, pos); // 提取name
        const argVal = nameValueArr[i].substring(pos + 1); // 提取value
        query[argName] = argVal;
      }
    }

    router.replace({ path, query }); // 要跳转的目标页面没有地址参数，直接跳，并且丢弃本登录页面历史路由
  });

}

// resetForm
const resetForm = () => {
  formData.value = { username: '', password: '' };
};

// 定义 formRef（校验规则）
const loading = ref(false);
const form = ref();
const formData = ref({ username: "", password: "" });
const rules = reactive({
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }]
});

onMounted(async () => {

  // 免登陆时直接跳转
  await enableLogin();
  if (tokenStore.needLogin == 'false') {
    router.replace({ path: '/' });  // 不需要登录
    return;
  }

  // 监听enter事件（调用登录）
  document.onkeydown = e => {
    e = window.event || e;
    if (e.code === "Enter" || e.code === "enter" || e.code === "NumpadEnter") {
      if (loading.value) return;
      login();
    }
  };
});
</script>

<style scoped lang="scss">
.login-container {
  position: relative;
  min-width: 520px;
  height: 100%;
  min-height: 520px;
  background-color: #eee;
  background-image: url('/image/loginbg.svg');
  background-size: 100% 100%;
  background-size: cover;

  .dark {
    position: absolute;
    top: 4.5%;
    right: 3.2%;
  }

  .login-box {
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: space-around;
    width: 96%;
    height: 94%;
    padding: 0 50px;
    background-color: hsl(0deg 0% 100% / 80%);
    border-radius: 10px;

    .login-left {
      width: 800px;

      img {
        width: 100%;
        height: 100%;
      }
    }

    .login-form {
      width: 420px;
      padding: 50px 40px 45px;
      background-color: #fff;
      border-radius: 10px;
      box-shadow: 2px 3px 7px rgb(0 0 0 / 20%);

      .login-logo {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 45px;

        .login-icon {
          width: 60px;
          height: 52px;
        }

        .logo-text {
          padding: 0 0 0 25px;
          margin: 0;
          font-size: 42px;
          font-weight: bold;
          color: #0081dd;
          white-space: nowrap;
        }
      }

      .el-form-item {
        margin-bottom: 40px;
      }

      .login-btn {
        display: flex;
        justify-content: space-between;
        width: 100%;
        margin-top: 40px;
        white-space: nowrap;

        .el-button {
          width: 185px;
        }
      }
    }
  }
}

.login-footer {
  position: absolute;
  bottom: 0;
  left: 50%;
  color: lightgrey;
  transform: translateX(-50%);
}

@media screen and (max-width: 1100px) {
  .login-left {
    display: none;
  }
}
</style>
