import { TableCell, TableRow } from "@mui/material";

//id [cid], name [first_name 'nick' second_name], email [email]
const UserDeletion = ({ deletion }) => (
  <>
    <TableRow>
      <TableCell>
        <div className="mono-bold removed">- {deletion.cid}</div>
      </TableCell>
      <TableCell>
        <div className="mono-bold removed">
          - {deletion.first_name} '{deletion.nick}' {deletion.second_name}
        </div>
      </TableCell>
      <TableCell>
        <div className="mono-bold removed">- {deletion.mail}</div>
      </TableCell>
    </TableRow>
  </>
);

export default UserDeletion;
