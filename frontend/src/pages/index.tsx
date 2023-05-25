import { Box } from "@mui/material";
import { useState } from "react";
import { SuggestionsTable } from "../components/SuggestionsTable";
import { GetServerSideProps, InferGetServerSidePropsType } from "next";
import { createGoldappsServerClient } from "../lib/goldapps/client";
import { checkLogin } from "../lib/goldapps/auth";
import { getSuggestions } from "../lib/goldapps/get-suggestions";
import { commitSuggestions } from "../lib/goldapps/commit-suggestions";
import { Suggestion, SuggestionWithState } from "../lib/goldapps/types";

export default function IndexPage({
  suggestions: serverSuggestions,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const [suggestions, setSuggestions] = useState<SuggestionWithState[]>(
    serverSuggestions.map(suggestion => ({
      ...suggestion,
      selected: false,
      error: null,
    })),
  );

  const onCommit = async () => {
    // Sending each change as a separate request allows us to show on which ones
    // the errors occur
    const suggestionsAfterCommit = await Promise.all(
      suggestions.map(async suggestion => {
        if (!suggestion.selected) {
          return suggestion;
        }

        try {
          await commitSuggestions([suggestion]);
          return suggestion;
        } catch (e) {
          return {
            ...suggestion,
            error: "Something went wrong, check the logs",
          };
        }
      }),
    );
    const remainingSuggestions = suggestionsAfterCommit.filter(
      suggestion => !suggestion.selected || suggestion.error,
    );
    setSuggestions(remainingSuggestions);
  };

  return (
    <Box sx={{ width: "100%" }}>
      {suggestions !== null && (
        <SuggestionsTable
          suggestions={suggestions}
          onSelectToggle={index => {
            if (index < 0 || index >= suggestions.length) {
              return;
            }

            setSuggestions(suggestions =>
              suggestions.map((suggestion, i) => {
                if (i === index) {
                  return {
                    ...suggestion,
                    selected: !suggestion.selected,
                  };
                } else {
                  return suggestion;
                }
              }),
            );
          }}
          onSelectAllToggle={() => {
            const allSelected = suggestions.every(
              suggestion => suggestion.selected,
            );

            setSuggestions(suggestions =>
              suggestions.map(suggestion => ({
                ...suggestion,
                selected: !allSelected,
              })),
            );
          }}
          onCommit={onCommit}
        />
      )}
    </Box>
  );
}

export const getServerSideProps: GetServerSideProps<{
  suggestions: Suggestion[];
}> = async ({ req }) => {
  const goldappsSessionCookie = req.cookies["goldapps-session"] || "";
  const goldappsClient = createGoldappsServerClient(goldappsSessionCookie);

  const redirectUri = await checkLogin(goldappsClient);
  if (redirectUri) {
    return {
      redirect: {
        destination: redirectUri,
        permanent: false,
      },
    };
  }

  const suggestions = await getSuggestions(goldappsClient);
  return {
    props: {
      suggestions,
    },
  };
};
