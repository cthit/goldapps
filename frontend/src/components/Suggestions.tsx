"use client";

import { useState } from "react";
import { Suggestion, SuggestionWithState } from "../lib/goldapps/types";
import { commitSuggestions } from "../lib/goldapps/commit-suggestions";
import { SuggestionsTable } from "./SuggestionsTable";

interface Props {
  suggestions: Suggestion[];
}

export const Suggestions = ({ suggestions: serverSuggestions }: Props) => {
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
    <div className="max-xl mx-2 pt-1 lg:container lg:mx-auto">
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
    </div>
  );
};
