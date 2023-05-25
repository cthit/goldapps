import { TableCell } from "@mui/material";
import { Group } from "../lib/goldapps/types";
import { getIdFromEmail } from "../utils";

interface Props {
  group: Group;
}

export const GroupAdditionRow = ({ group }: Props) => (
  <>
    <TableCell>
      <div className="mono-bold added">+ {getIdFromEmail(group.email)}</div>
    </TableCell>
    <TableCell>
      <div className="mono-bold added">+ {group.email}</div>
    </TableCell>
    <TableCell>
      {group.members &&
        group.members.map(group => (
          <div key={group} className="mono-bold added">
            + {group}
          </div>
        ))}
    </TableCell>
    <TableCell>Group Addition</TableCell>
  </>
);
