<template>
  <GxToolbar :height="themeStore.headerHeight">
    <template #left>
      <div style="display:flex;align-items:center;margin-left:8px;line-height: 30px;">
        <div :title="verInfo" style="width: 34px;color:white;text-align:center;cursor:pointer;" @click="clickLogo">
          <img src="/image/glc.png" style="width:34px;margin-top:9px;" />
        </div>
        <div style="text-align:center;">
          <img src="/image/title.png" style="height:32px;margin-top:11px;margin-left:26px" />
        </div>
      </div>
    </template>
    <template #right>
      <div v-if="needLogin" class="icon-item hover" style="margin-right:13px;">
        <el-dropdown style="margin-right:10px;" class="icon-item hover" trigger="click" size="default">
          <div class="item-user" style="min-width:90px;">
            <span style="margin-left:10px">欢迎您 {{ tokenStore.loginUserName }}</span>
            <svg-icon name="user" size="16" style="margin-right:0;margin-left:5px;" />
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="fnOpenChangePswDialog">
                <svg-icon name="edit-password" size="14" style="margin-right:8px;margin-left:0;" />修改密码
              </el-dropdown-item>
              <el-dropdown-item divided @click="logout">
                <svg-icon name="exit" size="14" style="margin-right:8px;margin-left:0;" />退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

      </div>
    </template>
  </GxToolbar>

  <!-- 修改密码 -->
  <GxDialog ref="dialog" title="修改密码" width="600" height="300" show-footer>
    <template #default>
      <div style="padding:0 60px;">
        <el-form ref="form" :model="formData" :rules="rules" label-width="auto" label-position="right" status-icon>
          <el-card shadow="never">

            <el-col :span="24">
              <el-form-item label="原密码" prop="oldPassword">
                <el-input v-model="formData.oldPassword" type="password" placeholder="请输入原密码"
                  autocomplete="new-password" />
              </el-form-item>
            </el-col>

            <el-form-item label="新密码" prop="password1">
              <el-input v-model="formData.password1" type="password" placeholder="请输入密码" autocomplete="new-password" />
            </el-form-item>

            <el-form-item label="确认密码" prop="password2">
              <el-input v-model="formData.password2" type="password" placeholder="请确认密码" autocomplete="new-password" />
            </el-form-item>

          </el-card>
        </el-form>
      </div>
    </template>
    <template #footer>
      <el-row style="justify-content: center;">
        <GxButton icon="select" type="primary" @click="fnChangePsw(form)">确 定</GxButton>
      </el-row>
    </template>
  </GxDialog>
</template>

<script setup>
import { gxUtil, useThemeStore, useTokenStore } from '~/pkgs';
import { userLogout } from '~/api';

const router = useRouter();
const verInfo = ref('');
const tokenStore = useTokenStore();
const themeStore = useThemeStore();
const formData = ref({ oldPassword: '', password1: '', password2: '' });
const form = ref(); // 表单实例
const dialog = ref();
const checkVersionDone = ref(false)

const validatePass2 = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请确认密码'))
  } else if (value !== formData.value.password1) {
    callback(new Error("两次输入密码不一致"))
  } else {
    callback()
  }
}
// 校验规则
const rules = reactive({
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' },],
  password1: [{ required: true, message: '请输入新密码', trigger: 'blur' },],
  password2: [{ required: true, validator: validatePass2, trigger: 'blur' },],
});

const headerHeight = computed(() => `${themeStore.headerHeight}px`);
const needLogin = computed(() => tokenStore.needLogin == 'true');
const headerColor = computed(() => {
  if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
    themeStore.headerColor = '#606266';
  }
  return themeStore.headerColor
});
const headerActiveColor = computed(() => {
  if (!themeStore.customHeaderColor) {
    if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
      themeStore.headerActiveColor = '#606266';
    } else {
      themeStore.headerActiveColor = themeStore.headerColor;
    }
  }
  return themeStore.headerActiveColor
});
const headerActiveBgColor = computed(() => {
  if (!themeStore.customHeaderColor) {
    if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
      themeStore.headerActiveBgColor = gxUtil.darkColor(themeStore.headerBgColor, 0.1);
    } else {
      themeStore.headerActiveBgColor = gxUtil.lightColor(themeStore.headerBgColor);
    }
  }
  return themeStore.headerActiveBgColor;
});
const svgLogoColor = computed(() => {
  if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
    return '#0081dc';
  }
  return '#eeeeee';
});

onMounted(() => checkVersion());

const fnChangePsw = (form) => {
  form.validate(valid => {
    if (!valid) {
      return;
    }

    const { oldPassword } = formData.value;
    const newPassword = formData.value.password1;
    const username = tokenStore.loginUserName

    $post('/v1/sysuser/changePsw', { username, oldPassword, newPassword }).then(rs => {
      if (!rs.success) {
        return $msg.error(rs.message);
      }
      $msg.info('操作成功！');
      dialog.value.show(false);
    });
  });
};

const fnOpenChangePswDialog = () => {
  formData.value.oldPassword = '';
  formData.value.password1 = '';
  formData.value.password2 = '';
  dialog.value.show(true);
};

async function logout() {
  if (await $msg.confirm('确定要退出系统吗？')) {
    userLogout();
    router.push('/glc/login');
  }
}

const clickLogo = () => {
  window.open('https://github.com/gotoeasy/glogcenter', '_blank');
};

function checkVersion() {
  if (!checkVersionDone.value) {
    checkVersionDone.value = true;
    // 从后台服务读取当前运行版本，避免多处维护版本号
    $post('/v1/version/info', {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
      if (rs.success) {
        verInfo.value = rs.result.version
        if (rs.result.latest && normalizeVer(rs.result.version) < normalizeVer(rs.result.latest)) {
          verInfo.value = `当前版本 ${rs.result.version} ，有新版本 ${rs.result.latest} 可更新`
        }
      }
    });
  }
}

// 0.1.2 => v100.1001.1002
function normalizeVer(ver) {
  const ary = ver.replace("v", "").split(".")
  return `v${100 + (ary[0] - 0)}.${1000 + (ary[1] - 0)}.${1000 + (ary[2] - 0)}`
}

</script>

<style lang="scss">
div.el-popper.el-dropdown__popper.el-popper {
  margin-top: -8px;
  margin-left: 8px;
}

.svg-logo {
  // height: v-bind('svgLogoHeight');
  fill: v-bind('svgLogoColor');
}
</style>

<style scoped lang="scss">
.icon-item {
  color: v-bind('headerColor');

}

.icon-item.hover {
  height: v-bind('headerHeight');
  line-height: v-bind('headerHeight');

  .svg-icon {
    margin: -2px 14px;
  }

  &:hover {
    color: v-bind('headerActiveColor');
    background-color: v-bind('headerActiveBgColor');

  }
}
</style>
