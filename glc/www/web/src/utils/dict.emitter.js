import { $emitter, useDictStore } from '~/pkgs';

let dictStore = null;

$emitter.on('$dict', (dictCode, research, code, editType) => {
  if (/^select/i.test(editType)) return fnDictOptions(dictCode, research); // 查下拉选项
  return code || code === 0 ? fnDictCodeName(dictCode, code) : ''; // 查名称
});

function fnDictOptions(dict, bSearch) {
  const dictCode = dict.replace('$dict:', '');

  if (!dictStore) {
    dictStore = useDictStore();
  }
  return dictStore.dict(dictCode, bSearch);
}

const map = new Map();
function fnDictCodeName(dict, code) {
  if (code == null || code == '') return '';

  dict = dict.replace('$dict:', '');
  const name = map.get(`${dict}#${code}`);
  if (!name) {
    const datas = fnDictOptions(dict) || [];
    for (let i = 0; i < datas.length; i++) {
      map.set(`${dict}#${datas[i].code}`, datas[i].name);
    }
  }
  return map.get(`${dict}#${code}`) || code;
}

export const optionData = fnDictOptions;
export const optionName = fnDictCodeName;
