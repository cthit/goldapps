import { TableCell } from "@mui/material";
import { User } from "../lib/goldapps/types";

interface Props {
  user: User;
}

export const UserAdditionRow = ({ user }: Props) => (
  <>
    <TableCell>
      <div className="mono-bold added">+ {user.cid}</div>
    </TableCell>
    <TableCell>
      <div className="mono-bold added">
        + {user.first_name} '{user.nick}' {user.second_name}
      </div>
    </TableCell>
    <TableCell>
      <div className="mono-bold added">+ {user.mail}</div>
    </TableCell>
    <TableCell>User Addition</TableCell>
  </>
);
