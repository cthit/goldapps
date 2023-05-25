import {
  Box,
  Button,
  Checkbox,
  Paper,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import React from "react";
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
    <Box sx={{ padding: "2rem" }}>
      <TableContainer component={Paper}>
        <TableHead>
          <TableRow>
            <TableCell padding="checkbox">
              <Checkbox checked={allSelected} onChange={onSelectAllToggle} />
            </TableCell>
            <TableCell>Id</TableCell>
            <TableCell>Name</TableCell>
            <TableCell>E-mail(s)</TableCell>
            <TableCell>Type</TableCell>
            <TableCell>Error</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {suggestions.map((suggestion, suggestionIndex) => (
            <TableRow key={suggestionIndex}>
              <TableCell padding="checkbox">
                <Checkbox
                  checked={suggestion.selected}
                  onChange={() => onSelectToggle(suggestionIndex)}
                />
              </TableCell>
              {transformSuggestionToRow(suggestion)}
              <TableCell>{suggestion.error}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </TableContainer>
      <div style={{ display: "flex", justifyContent: "flex-end" }}>
        <Button
          sx={{ marginTop: "2rem" }}
          variant="outlined"
          onClick={onCommit}
        >
          Commit
        </Button>
      </div>
    </Box>
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
