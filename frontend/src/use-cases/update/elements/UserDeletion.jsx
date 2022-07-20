import { Checkbox, TableCell, TableRow } from "@mui/material";

//id [cid], name [first_name 'nick' second_name], email [email]
const UserDeletion = ({ deletion, selected, onChange }) => (
  <>
    <TableRow>
      <TableCell padding="checkbox">
        <Checkbox
          checked={selected.includes(deletion.id)}
          onChange={() => onChange(deletion.id)}
        />
      </TableCell>
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
      <TableCell>User Deletion</TableCell>
    </TableRow>
  </>
);

export default UserDeletion;
