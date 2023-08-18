import { axios } from './axios.config';
import { theme } from './theme.config';

export const config = {
  ...axios,
  ...theme,
};
