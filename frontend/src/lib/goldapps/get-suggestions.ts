import { AxiosInstance } from "axios";
import { Suggestions } from "./types";
import { flattenSuggestions } from "./transform";

export async function getSuggestions(client: AxiosInstance) {
  const { data: suggestions } = await client.get<Suggestions>(
    "/api/suggestions",
    {
      params: {
        from: "gamma",
        to: "gapps",
      },
    },
  );

  return flattenSuggestions(suggestions);
}
