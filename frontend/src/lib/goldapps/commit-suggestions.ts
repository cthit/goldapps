import { createGoldappsBrowserClient } from "./client-browser";
import { unflattedSuggestions } from "./transform";
import { Suggestion } from "./types";

export async function commitSuggestions(flatSuggestions: Suggestion[]) {
  const client = createGoldappsBrowserClient();
  const suggestions = unflattedSuggestions(flatSuggestions);
  await client.post("/api/commit", suggestions, {
    params: {
      to: "gapps",
    },
  });
}
