import { Suggestion, Suggestions } from "./types";

export function flattenSuggestions(suggestions: Suggestions) {
  let flatSuggestions: Suggestion[] = [];
  const { userChanges, groupChanges } = suggestions;
  if (userChanges) {
    const { additions, deletions, userUpdates } = userChanges;
    if (additions) {
      for (const user of additions) {
        flatSuggestions.push({
          type: "AddUser",
          user,
        });
      }
    }

    if (deletions) {
      for (const user of deletions) {
        flatSuggestions.push({
          type: "DeleteUser",
          user,
        });
      }
    }

    if (userUpdates) {
      for (const { before, after } of userUpdates) {
        flatSuggestions.push({
          type: "ChangeUser",
          before,
          after,
        });
      }
    }
  }

  if (groupChanges) {
    const { additions, deletions, groupUpdates } = groupChanges;
    if (additions) {
      for (const group of additions) {
        flatSuggestions.push({
          type: "AddGroup",
          group,
        });
      }
    }

    if (deletions) {
      for (const group of deletions) {
        flatSuggestions.push({
          type: "DeleteGroup",
          group,
        });
      }
    }

    if (groupUpdates) {
      for (const { before, after } of groupUpdates) {
        flatSuggestions.push({
          type: "ChangeGroup",
          before,
          after,
        });
      }
    }
  }

  return flatSuggestions;
}

export function unflattedSuggestions(flatSuggestions: Suggestion[]) {
  const suggestions: Suggestions = {
    userChanges: {
      additions: [],
      deletions: [],
      userUpdates: [],
    },
    groupChanges: {
      additions: [],
      deletions: [],
      groupUpdates: [],
    },
  };

  for (const suggestion of flatSuggestions) {
    switch (suggestion.type) {
      case "AddUser":
        suggestions.userChanges!.additions!.push(suggestion.user);
        break;
      case "DeleteUser":
        suggestions.userChanges!.deletions!.push(suggestion.user);
        break;
      case "ChangeUser":
        suggestions.userChanges!.userUpdates!.push({
          before: suggestion.before,
          after: suggestion.after,
        });
        break;
      case "AddGroup":
        suggestions.groupChanges!.additions!.push(suggestion.group);
        break;
      case "DeleteGroup":
        suggestions.groupChanges!.deletions!.push(suggestion.group);
        break;
      case "ChangeGroup":
        suggestions.groupChanges!.groupUpdates!.push({
          before: suggestion.before,
          after: suggestion.after,
        });
        break;
    }
  }

  return suggestions;
}
