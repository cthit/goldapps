import { Suggestion } from "../lib/goldapps/types";
import { Suggestions } from "../components/Suggestions";
import { cookies } from "next/headers";
import { createGoldappsServerClient } from "../lib/goldapps/client";
import { checkLogin } from "../lib/goldapps/auth";
import { redirect } from "next/navigation";
import { getSuggestions } from "../lib/goldapps/get-suggestions";

async function fetchSuggestions() {
  const cookieJar = cookies();
  const goldappsSessionCookie = cookieJar.get("goldapps-session")?.value || "";
  const goldappsClient = createGoldappsServerClient(goldappsSessionCookie);
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
