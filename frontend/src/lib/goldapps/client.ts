import axios from "axios";

export function createGoldappsServerClient(goldappsSession: string) {
  return axios.create({
    headers: {
      Cookie: `goldapps-session=${goldappsSession}`,
    },
    // TODO: get from config
    baseURL: "http://localhost:8080",
  });
}

export function createGoldappsBrowserClient() {
  return axios.create({
    withCredentials: true,
    baseURL: "http://localhost:8080",
  });
}
