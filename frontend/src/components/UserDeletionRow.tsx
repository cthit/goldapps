import { TableCell } from "@mui/material";
import { User } from "../lib/goldapps/types";

interface Props {
  user: User;
}

export const UserDeletionRow = ({ user }: Props) => (
  <>
    <TableCell>
      <div className="mono-bold removed">- {user.cid}</div>
    </TableCell>
    <TableCell>
      <div className="mono-bold removed">
        - {user.first_name} '{user.nick}' {user.second_name}
      </div>
    </TableCell>
    <TableCell>
      <div className="mono-bold removed">- {user.mail}</div>
    </TableCell>
    <TableCell>User Deletion</TableCell>
  </>
);
