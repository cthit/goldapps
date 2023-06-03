import axios from "axios";

export function createGoldappsBrowserClient() {
  return axios.create({
    withCredentials: true,
  });
}
