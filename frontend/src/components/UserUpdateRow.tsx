import { TableCell } from "@mui/material";
import { User } from "../lib/goldapps/types";

interface Props {
  before: User;
  after: User;
}

export const UserUpdateRow = ({ before, after }: Props) => (
  <>
    <TableCell>{before.cid}</TableCell>
    <TableCell>
      {!userNameEqual(before, after) ? (
        <>
          <div className="mono-bold removed">
            - {before.first_name} '{before.nick}' {before.second_name}
          </div>
          <div className="mono-bold added">
            + {after.first_name} '{after.nick}' {after.second_name}
          </div>
        </>
      ) : (
        <>
          {before.first_name} '{before.nick}' {before.second_name}
        </>
      )}
    </TableCell>
    <TableCell>
      {!userEmailEqual(before, after) ? (
        <>
          <div className="mono-bold removed">- {before.mail}</div>
          <div className="mono-bold added">+ {after.mail}</div>
        </>
      ) : (
        <>{before.mail}</>
      )}
    </TableCell>
    <TableCell>User Update</TableCell>
  </>
);

function userNameEqual(before: User, after: User) {
  return (
    before.first_name === after.first_name &&
    before.nick === after.nick &&
    before.second_name === after.second_name
  );
}

function userEmailEqual(before: User, after: User) {
  return before.mail === after.mail;
}
