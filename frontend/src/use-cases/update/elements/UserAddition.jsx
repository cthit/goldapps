import { Checkbox, TableCell, TableRow } from "@mui/material";

//id [cid], name [first_name 'nick' second_name], email [email]
const UserAddition = ({ addition, selected, onChange }) => (
  <>
    <TableRow>
      <TableCell padding="checkbox">
        <Checkbox
          checked={selected.includes(addition.id)}
          onChange={() => onChange(addition.id)}
        />
      </TableCell>
      <TableCell>
        <div className="mono-bold added">+ {addition.cid}</div>
      </TableCell>
      <TableCell>
        <div className="mono-bold added">
          + {addition.first_name} '{addition.nick}' {addition.second_name}
        </div>
      </TableCell>
      <TableCell>
        <div className="mono-bold added">+ {addition.mail}</div>
      </TableCell>
    </TableRow>
  </>
);

export default UserAddition;
