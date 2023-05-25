import { TableCell } from "@mui/material";
import { Group } from "../lib/goldapps/types";
import { getIdFromEmail } from "../utils";

interface Props {
  group: Group;
}

export const GroupDeletionRow = ({ group }: Props) => (
  <>
    <TableCell>
      <div className="mono-bold removed">- {getIdFromEmail(group.email)}</div>
    </TableCell>
    <TableCell>
      <div className="mono-bold removed">- {group.email}</div>
    </TableCell>
    <TableCell>
      {group.members &&
        group.members.map(group => (
          <div key={group} className="mono-bold removed">
            - {group}
          </div>
        ))}
    </TableCell>
    <TableCell>Group Deletion</TableCell>
  </>
);
