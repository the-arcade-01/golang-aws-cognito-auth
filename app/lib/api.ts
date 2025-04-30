import axios from "axios";
import AsyncStorage from "@react-native-async-storage/async-storage";

const API_URL = process.env.EXPO_PUBLIC_API_URL;

const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const signup = async (name: string, email: string, password: string) => {
  const response = await api.post("/signup", { name, email, password });
  return response.data;
};

export const confirmEmail = async (email: string, code: string) => {
  const response = await api.post("/confirm", { email, code });
  return response.data;
};

export const login = async (email: string, password: string) => {
  const response = await api.post("/login", { email, password });
  if (response.data.access_token) {
    await AsyncStorage.setItem("token", response.data.access_token);
  }
  return response.data;
};

export const getUserInfo = async () => {
  const response = await api.get("/user/info");
  return response.data;
};

export const logout = async () => {
  await AsyncStorage.removeItem("token");
};
