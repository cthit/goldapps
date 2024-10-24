import { Suggestions } from "../components/Suggestions";
import { checkLogin } from "../lib/goldapps/auth";
import { redirect } from "next/navigation";
import { getSuggestions } from "../lib/goldapps/get-suggestions";
import { createGoldappsServerClient } from "../lib/goldapps/client-server";



export default async function IndexPage() {
  const fetchSuggestions = async () => {
    const goldappsClient = createGoldappsServerClient();
    return getSuggestions(goldappsClient);
  }

  const suggestions = await fetchSuggestions();
  return <Suggestions suggestions={suggestions} />;
}
