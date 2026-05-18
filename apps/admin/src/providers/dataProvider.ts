import dataProviderSimpleRest from "@refinedev/simple-rest";
import axios from "axios";

const API_URL = "http://localhost:3001/api/admin";

// Configure axios interceptor to add JWT token
axios.interceptors.request.use((config) => {
  const auth = localStorage.getItem("auth");
  if (auth) {
    const { access_token } = JSON.parse(auth);
    config.headers.Authorization = `Bearer ${access_token}`;
  }
  return config;
});

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("auth");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export const dataProvider = dataProviderSimpleRest(API_URL, axios);
