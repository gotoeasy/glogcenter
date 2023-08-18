import ContextMenu from '@imengyu/vue3-context-menu';

export const openConfigContextmenu = ({ event, label = '设定...　　　', icon = 'menu', click = () => {} }) => {
  event.preventDefault();
  ContextMenu.showContextMenu({
    x: event.x + 1,
    y: event.y,
    items: [
      {
        label,
        svgIcon: `#gx_icon-${icon}`,
        onClick: click,
      },
    ],
  });
};
