'use server';

import { AxiosResponse } from "axios";
import { cookies } from "next/headers";
import { createGoldappsServerClient } from "./client-server";

// Returns login uri if not logged in
export async function checkLogin() {
  try {
    const client = createGoldappsServerClient();
    const { status: checkLoginStatus, data, headers } = await client.get<string>(
      "/api/checkLogin",
    );
    if (checkLoginStatus !== 200) {
      const serverCookies = headers["set-cookie"]?.map((cookie) => {
        const [nameValue, ...rest] = cookie.split(';');
        const [name, value] = nameValue.split('=');
        return { name, value, attributes: rest.join(';') };
      }) || [];

      const oauthStateCookie = serverCookies.find((cookie) => cookie.name === "oauth_state");

      if (oauthStateCookie) {
        cookies().set({
          name: "oauth_state",
          value: oauthStateCookie.value,
          maxAge: 3600, // 1 hour
          httpOnly: true,
          secure: true,
          sameSite: "none"
        });
      }
      
      return data;
    }

    return null;
  } catch (e) {
    const response = (e as any).response as AxiosResponse<string>;
    return response.data;
  }
}
