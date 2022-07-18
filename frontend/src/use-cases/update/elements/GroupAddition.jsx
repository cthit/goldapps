import { TableCell, TableRow } from "@mui/material";
import { getId } from "../../../utils/utils";

//id [<name before @chalmers.it>], name [email], email [member emails]
const GroupAddition = ({ addition }) => (
  <>
    <TableRow>
      <TableCell>
        <div className="mono-bold added">+ {getId(addition.email)}</div>
      </TableCell>
      <TableCell>
        <div className="mono-bold added">+ {addition.email}</div>
      </TableCell>
      <TableCell>
        {addition.members.map(m => (
          <div className="mono-bold added">+ {m}</div>
        ))}
      </TableCell>
    </TableRow>
  </>
);

export default GroupAddition;