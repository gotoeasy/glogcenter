const { defineConfig } = require('eslint-define-config');

module.exports = defineConfig({
  root: true,
  env: {
    browser: true,
    node: true,
    jest: true,
    es6: true,
  },
  plugins: ['vue'],
  parser: 'vue-eslint-parser',
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    allowImportExportEverywhere: true,
    ecmaFeatures: {
      jsx: true,
    },
  },
  extends: [
    'eslint-config-airbnb-base',
    'eslint:recommended',
    'plugin:vue/vue3-essential',
    'plugin:vue/vue3-recommended',
    'prettier',
    // 'plugin:prettier/recommended',
    './.eslintrc-auto-import.json',
  ],
  rules: {
    'import/no-extraneous-dependencies': 0, // 禁止使用多余的包
    'import/extensions': 0, // 确保在导入路径内一致使用文件扩展名
    'import/no-unresolved': 0, // 确保导入指向可以解析的文件/模块
    'import/prefer-default-export': 0, // 首选默认导出导入/首选默认导出
    'no-var': 'error', // 要求使用 let 或 const 而不是 var
    'no-new': 1, // 禁止使用 new 以避免产生副作用
    'no-shadow': 0, // 禁止变量声明与外层作用域的变量同名
    'no-console': 0, // 禁用 console
    'no-underscore-dangle': 0, // 禁止标识符中有悬空下划线
    'no-confusing-arrow': 0, // 禁止在可能与比较操作符相混淆的地方使用箭头函数
    'no-plusplus': 0, // 禁用一元操作符 ++ 和 --
    'no-param-reassign': 0, // 禁止对 function 的参数进行重新赋值
    'no-restricted-syntax': 0, // 禁用特定的语法
    'no-use-before-define': 0, // 禁止在变量定义之前使用它们
    'no-prototype-builtins': 0, // 禁止直接调用 Object.prototypes 的内置属性
    'no-unneeded-ternary': 'error', // 禁止可以在有更简单的可替代的表达式时使用三元操作符
    'no-duplicate-imports': 'error', // 禁止重复模块导入
    'no-useless-computed-key': 'error', // 禁止在对象中使用不必要的计算属性
    'no-useless-escape': 0, // 禁止不必要的转义字符
    'no-continue': 0, // 禁用 continue 语句
    indent: ['error', 2, { SwitchCase: 1 }], // 强制使用一致的缩进
    camelcase: 0, // 强制使用骆驼拼写法命名约定
    'class-methods-use-this': 0, // 强制类方法使用 this
    'new-cap': 0, // 要求构造函数首字母大写
    'func-style': 0, // 强制一致地使用 function 声明或表达式
    'max-len': 0, // 强制一行的最大长度
    'consistent-return': 0, // 要求 return 语句要么总是指定返回的值，要么不指定
    'default-case': 2, // 强制switch要有default分支
    'rest-spread-spacing': 'error', // 强制剩余和扩展运算符及其表达式之间有空格
    'prefer-const': 'error', // 要求使用 const 声明那些声明后不再被修改的变量
    'arrow-spacing': 'error', // 强制箭头函数的箭头前后使用一致的空格
    'prefer-destructuring': ['error', { object: true, array: false }], // 只强制对象解构，不强制数组解构
    'no-multiple-empty-lines': [
      1,
      {
        max: 1, // 空行最多不能超过1行(多个空行则合并)
      },
    ],
    'vue/multi-word-component-names': 'off',
    'no-debugger': 0, // 使用debugger不警告
    'no-unused-expressions': 0, // 禁止出现未使用过的表达式(如：fn && fn())
    'no-return-assign': 0, // 禁止在 return 语句中使用赋值语句
    'vue/first-attribute-linebreak': 0, // 忽略换行检查
    'no-bitwise': 0, // 允许位运算
    'vue/valid-v-slot': 0,
    eqeqeq: 0,
    'vue/no-v-html': 'off',
  },
});
