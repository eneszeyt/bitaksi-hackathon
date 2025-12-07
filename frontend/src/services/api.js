import axios from 'axios';

// api gateway address
const api = axios.create({
  baseURL: 'http://localhost:8000', 
});

// interceptor configuration to automatically add token to headers
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// login function
// backend expects form data (x-www-form-urlencoded), not json
// so we use URLSearchParams to format the payload correctly
export const login = (username, password) => {
  const params = new URLSearchParams();
  params.append('username', username);
  params.append('password', password);
  
  return api.post('/login', params);
};

// get list of all drivers with pagination
export const getDrivers = () => api.get('/drivers?page=1&pageSize=100');

// get nearby drivers based on coordinates
// type parameter is optional (e.g., 'yellow', 'black' or empty for all)
export const getNearbyDrivers = (lat, lon, type = '') => {
    let url = `/drivers/nearby?lat=${lat}&lon=${lon}`;
    if (type) {
        url += `&taxiType=${type}`;
    }
    return api.get(url);
};

export default api;