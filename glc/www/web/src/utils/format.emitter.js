// --------------------------------------------------
// 这个文件里定义各钟格式转换处理，用于表格的展示转换
// --------------------------------------------------
import { $emitter } from '~/pkgs';

// 金额或数字加逗号（3位1撇）
$emitter.on('$format:amount', value => (value == null ? value : formatNumber(Number(value))));

function formatNumber(num) {
  // 将数字转换为字符串
  const str = String(num);

  // 检查是否为负数或小数
  const isNegative = str.includes('-');
  const isDecimal = str.includes('.');

  // 分割整数部分和小数部分
  const parts = str.split('.');
  let integerPart = parts[0];
  const decimalPart = parts[1] || '';

  // 添加千分号
  integerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ',');

  // 重新组合整数和小数部分
  let result = (isNegative ? '-' : '') + integerPart;
  if (isDecimal) {
    result += `.${decimalPart}`;
  }

  return result;
}

export const formatter = (fmtKey, value, item) => {
  if (!fmtKey) {
    console.warn(`参数有误，fmtKey无效`);
    return value;
  }

  const fn = $emitter.on(fmtKey);
  if (fn) return fn(value, item);

  console.warn(`找不到格式转换器${fmtKey}，直接返回不作处理`);
  return value;
};

// 方式
$emitter.on('$format', formatter);
