import { UserAdditionRow } from "./UserAdditionRow";
import { UserDeletionRow } from "./UserDeletionRow";
import { UserUpdateRow } from "./UserUpdateRow";
import { GroupAdditionRow } from "./GroupAdditionRow";
import { GroupDeletionRow } from "./GroupDeletionRow";
import { GroupUpdateRow } from "./GroupUpdateRow";
import { SuggestionWithState } from "../lib/goldapps/types";

interface Props {
  suggestions: SuggestionWithState[];
  onSelectToggle: (index: number) => void;
  onSelectAllToggle: () => void;
  onCommit: () => void;
}

export const SuggestionsTable = ({
  suggestions,
  onSelectToggle,
  onSelectAllToggle,
  onCommit,
}: Props) => {
  const allSelected = suggestions.every(suggestion => suggestion.selected);
  return (
    <>
      <table className="w-full border-collapse text-left">
        <thead>
          <tr>
            <th>
              <input
                type="checkbox"
                checked={allSelected}
                onChange={onSelectAllToggle}
              />
            </th>
            <th>Id</th>
            <th>Name</th>
            <th>E-mail(s)</th>
            <th>Type</th>
            <th>Error</th>
          </tr>
        </thead>
        <tbody>
          {suggestions.map((suggestion, suggestionIndex) => (
            <tr key={suggestionIndex} className="border-b border-t py-5">
              <td>
                <input
                  type="checkbox"
                  checked={suggestion.selected}
                  onChange={() => onSelectToggle(suggestionIndex)}
                />
              </td>
              {transformSuggestionToRow(suggestion)}
              <td>{suggestion.error}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className="flex justify-end">
        <button
          className="mt-2 rounded bg-blue-500 px-4 py-2 text-white"
          onClick={onCommit}
        >
          Commit
        </button>
      </div>
    </>
  );
};

function transformSuggestionToRow(suggestion: SuggestionWithState) {
  switch (suggestion.type) {
    case "AddUser":
      return <UserAdditionRow user={suggestion.user} />;
    case "DeleteUser":
      return <UserDeletionRow user={suggestion.user} />;
    case "ChangeUser":
      return (
        <UserUpdateRow before={suggestion.before} after={suggestion.after} />
      );
    case "AddGroup":
      return <GroupAdditionRow group={suggestion.group} />;
    case "DeleteGroup":
      return <GroupDeletionRow group={suggestion.group} />;
    case "ChangeGroup":
      return (
        <GroupUpdateRow before={suggestion.before} after={suggestion.after} />
      );
  }
}
