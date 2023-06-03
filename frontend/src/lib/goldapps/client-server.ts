import axios from "axios";
import { cookies } from "next/headers";

const GOLDAPPS_COOKIE_NAME = "goldapps-session";

export function createGoldappsServerClient() {
  const cookieJar = cookies();
  const goldappsCookie = cookieJar.get(GOLDAPPS_COOKIE_NAME)?.value || "";
  return axios.create({
    headers: {
      Cookie: `${GOLDAPPS_COOKIE_NAME}=${goldappsCookie}`,
    },
    baseURL: String(process.env.GOLDAPPS_URL),
  });
}
