'use server';

import { AxiosHeaders, AxiosResponse, RawAxiosResponseHeaders } from "axios";
import { cookies } from "next/headers";
import { createGoldappsServerClient } from "./client-server";
import { ResponseCookie } from "next/dist/compiled/@edge-runtime/cookies";

function extractOauthStateCookie(headers: RawAxiosResponseHeaders | (RawAxiosResponseHeaders & AxiosHeaders)): ResponseCookie | null {
  if (headers === null || headers === undefined) {
    return null;
  }
  const setCookieHeader = headers["set-cookie"] || null;
  const serverCookies = (Array.isArray(setCookieHeader) ? setCookieHeader : [setCookieHeader]).map((cookie) => {
    if (typeof cookie === "string") {
      const [nameValue, ...rest] = cookie.split(";");
      const [name, value] = nameValue.split("=");
      return { name, value, attributes: rest.join(";") };
    }
    return null;
  }).filter(cookie => cookie !== null) || [];

  const oauthStateCookie = serverCookies.find((cookie: any) => cookie.name === "oauth_state");

  if (oauthStateCookie) {
    return {
      name: "oauth_state",
      value: oauthStateCookie.value,
      maxAge: 3600, // 1 hour
      httpOnly: true,
      secure: true,
      sameSite: "none",
    };
  }
  return null
}

// Returns login uri if not logged in
export async function checkLogin() : Promise<{data: string, cookie: (ResponseCookie | null)} | null> {
  try {

    const client = createGoldappsServerClient();
    const { status: checkLoginStatus, data, headers } = await client.get<string>(
      "/api/checkLogin"
    );
    if (checkLoginStatus !== 200) {
      return {data, cookie: extractOauthStateCookie(headers)};
    }

    return null;
  } catch (e) {
    const response = (e as any).response as AxiosResponse<string>;
    return {data: response.data, cookie: extractOauthStateCookie(response.headers)};
  }
}
