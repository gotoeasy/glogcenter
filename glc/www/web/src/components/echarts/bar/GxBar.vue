<template>
  <div ref="div" style="width: 600px;height:400px;" @contextmenu="fnContextmenu"></div>
</template>

<script setup >
import * as echarts from 'echarts'
import { openConfigContextmenu } from '~/utils'

const div = ref()

const fnContextmenu = event => {
  openConfigContextmenu({
    event,
    click() {
      alert()
    }
  });
}

let myChart = null;

onMounted(() => {
  myChart = echarts.init(div.value, null, { renderer: 'svg' });

  // 指定图表的配置项和数据
  const option = {
    legend: {
      data: ['销量', '单价']
    },
    xAxis: {
      data: ['衬衫', '羊毛衫', '雪纺衫', '裤子', '高跟鞋', '袜子']
    },
    yAxis: {},
    series: [
      {
        name: '销量',
        type: 'bar',
        data: [5, 20, 36, 10, 10, 20]
      },
      {
        name: '单价',
        type: 'bar',
        data: [53.1, 34.3, 36, 40, 30, 20]
      }
    ]
  };

  myChart.setOption(option);  // 使用刚指定的配置项和数据显示图表。
  window.addEventListener('resize', () => {
    myChart.resize();  // 改变大小时重绘
  });
});

onUnmounted(() => {
  myChart.dispose(); // 销毁
});

</script>
