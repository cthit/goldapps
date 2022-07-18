import { Checkbox, TableCell, TableRow } from "@mui/material";
import { getId } from "../../../utils/utils";

//id [<name before @chalmers.it>], name [email], email [member emails]
const GroupDeletion = ({ deletion, selected, onChange }) => (
  <>
    <TableRow>
      <TableCell padding="checkbox">
        <Checkbox
          checked={selected.includes(deletion.id)}
          onChange={() => onChange(deletion.id)}
        />
      </TableCell>
      <TableCell>
        <div className="mono-bold removed">- {getId(deletion.email)}</div>
      </TableCell>
      <TableCell>
        <div className="mono-bold removed">- {deletion.email}</div>
      </TableCell>
      <TableCell>
        {deletion.members.map(m => (
          <div className="mono-bold removed">- {m}</div>
        ))}
      </TableCell>
    </TableRow>
  </>
);

export default GroupDeletion;
