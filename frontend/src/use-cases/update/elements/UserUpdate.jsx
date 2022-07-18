import { Checkbox, TableCell, TableRow } from "@mui/material";

//id [cid], name [first_name 'nick' second_name], email [email]
const UserUpdate = ({ change, selected, onChange }) => (
  <>
    <TableRow>
      <TableCell padding="checkbox">
        <Checkbox
          checked={selected.includes(change.id)}
          onChange={() => onChange(change.id)}
        />
      </TableCell>
      <TableCell>{change.before.cid}</TableCell>
      <TableCell>
        {change.before.first_name !== change.after.first_name ||
        change.before.nick !== change.after.nick ||
        change.before.second_name !== change.after.second_name ? (
          <>
            <div className="mono-bold removed">
              - {change.before.first_name} '{change.before.nick}'{" "}
              {change.before.second_name}
            </div>
            <div className="mono-bold added">
              + {change.after.first_name} '{change.after.nick}'{" "}
              {change.after.second_name}
            </div>
          </>
        ) : (
          <>
            {change.before.first_name} '{change.before.nick}'{" "}
            {change.before.second_name}
          </>
        )}
      </TableCell>
      <TableCell>
        {change.before.mail !== change.after.mail ? (
          <>
            <div className="mono-bold removed">- {change.before.mail}</div>
            <div className="mono-bold added">+ {change.after.mail}</div>
          </>
        ) : (
          <>{change.before.mail}</>
        )}
      </TableCell>
    </TableRow>
  </>
);

export default UserUpdate;
