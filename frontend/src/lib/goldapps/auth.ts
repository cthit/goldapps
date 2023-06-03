import { AxiosInstance, AxiosResponse } from "axios";

// Returns login uri if not logged in
export async function checkLogin(client: AxiosInstance) {
  try {
    const { status: checkLoginStatus, data } = await client.get<string>(
      "/api/checkLogin",
    );
    if (checkLoginStatus !== 200) {
      return data;
    }

    return null;
  } catch (e) {
    const response = (e as any).response as AxiosResponse<string>;
    return response.data;
  }
}
