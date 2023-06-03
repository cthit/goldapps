import { Suggestions } from "../components/Suggestions";
import { checkLogin } from "../lib/goldapps/auth";
import { redirect } from "next/navigation";
import { getSuggestions } from "../lib/goldapps/get-suggestions";
import { createGoldappsServerClient } from "../lib/goldapps/client-server";

async function fetchSuggestions() {
  const goldappsClient = createGoldappsServerClient();
  const redirectUri = await checkLogin(goldappsClient);
  if (redirectUri) {
    return redirect(redirectUri);
  }

  return getSuggestions(goldappsClient);
}

export default async function IndexPage() {
  const suggestions = await fetchSuggestions();
  return <Suggestions suggestions={suggestions} />;
}
