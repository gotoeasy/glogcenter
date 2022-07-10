import axios from 'axios'

const BASE_URL = process.env.NODE_ENV === 'production' ? location.origin : 'http://127.0.0.1:8080'

function request(url, method = 'get', data) {
  return axios({
    baseURL: BASE_URL,
    url: url,
    method: method,
    data: data,
  })
}

function post(url, formData) {
  formData.append('token', sessionStorage.getItem("glctoken"));
  return axios.post(url,formData,{
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    }
  })
}

export default {
  search(params={}) {

    let formData = new FormData();
    for(let k in params){
      formData.append(k, params[k]);
    }
    return post(`${BASE_URL}/glc/v1/log/search`, formData)
  },
  searchStorageNames() {
    return post(`${BASE_URL}/glc/v1/store/names`, new FormData())
  },
  searchStorages(params={}) {
    let formData = new FormData();
    for(let k in params){
      formData.append(k, params[k]);
    }
    return post(`${BASE_URL}/glc/v1/store/list`, formData)
  },
  deleteStorage(name) {
    let formData = new FormData();
    formData.append("storeName", name);
    return post(`${BASE_URL}/glc/v1/store/delete`, formData)
  },
  login(user,pass) {
    let formData = new FormData();
    formData.append("username", user);
    formData.append("password", pass);
    return post(`${BASE_URL}/glc/v1/user/login`, formData)
  },
  enableLogin() {
    let formData = new FormData();
    return post(`${BASE_URL}/glc/v1/user/enableLogin`, formData)
  },
}
