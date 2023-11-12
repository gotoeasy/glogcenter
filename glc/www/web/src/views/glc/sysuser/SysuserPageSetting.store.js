import { usePiniaStore } from '~/pkgs';

export default usePiniaStore({
  id: 'useSysuserPageSettingStore',
  state: () => ({
    detailMode: 'dialog', // panel、drawer、dialog
  }),
  getters: {
    pageOptions() {
      return {
        pageTitle: '用户',
        urlAdd: '/v1/sysuser/save',
        urlEdit: '/v1/sysuser/save',
        urlSearch: '/v1/sysuser/list',
        urlDelete: '/v1/sysuser/del',
      };
    },
  },
  storage: localStorage,
});
