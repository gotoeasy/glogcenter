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
  return axios.post(url,formData,{
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    //  'X-GLC-AUTH': 'glogcenter',
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
  // remove(db, id) {
  //   return request(`/remove?database=${db}`, 'post', { id })
  // },
  // gc() {
  //   return request('/gc')
  // },
  // getStatus() {
  //   return request('/status')
  // },
  // addIndex(db, index) {
  //   return request(`/index?database=${db}`, 'post', index )
  // },
  // drop(db){
  //   return request(`/db/drop?database=${db}`)
  // },
  // create(db){
  //   return request(`/db/create?database=${db}`)
  // }
}
