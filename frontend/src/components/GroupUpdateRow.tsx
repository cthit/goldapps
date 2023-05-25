import { TableCell } from "@mui/material";
import { useMemo } from "react";
import { Group } from "../lib/goldapps/types";
import { arrayDiff, getIdFromEmail } from "../utils";

interface Props {
  before: Group;
  after: Group;
}

export const GroupUpdateRow = ({ before, after }: Props) => {
  const { additions, deletions, unchanged } = useMemo(
    () => arrayDiff(before.members || [], after.members || []),
    [before.members, after.members],
  );

  return (
    <>
      <TableCell>{getIdFromEmail(before.email)}</TableCell>
      <TableCell>{before.email}</TableCell>
      <TableCell>
        {unchanged.map(member => (
          <div>{member}</div>
        ))}
        {deletions.map(member => (
          <div className="mono-bold removed">- {member}</div>
        ))}
        {additions.map(member => (
          <div className="mono-bold added">+ {member}</div>
        ))}
      </TableCell>
      <TableCell>Group Update</TableCell>
    </>
  );
};
