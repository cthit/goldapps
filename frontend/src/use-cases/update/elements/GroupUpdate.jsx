import { Checkbox, TableCell, TableRow } from "@mui/material";
import { getId } from "../../../utils/utils";

const getDiff = (oldArr, newArr) => {
  let res = [];
  for (let i of oldArr) {
    if (newArr.includes(i)) {
      res.push(<div>{i}</div>);
    }
  }
  for (let i of oldArr) {
    if (!newArr.includes(i)) {
      res.push(<div className="mono removed">- {i}</div>);
    }
  }
  for (let i of newArr) {
    if (!oldArr.includes(i)) {
      res.push(<div className="mono added">+ {i}</div>);
    }
  }
  return res;
};

//id [<name before @chalmers.it>], name [email], email [member emails]
const GroupUpdate = ({ change, selected, onChange }) => (
  <>
    <TableRow>
      <TableCell padding="checkbox">
        <Checkbox
          checked={selected.includes(change.id)}
          onChange={() => onChange(change.id)}
        />
      </TableCell>
      <TableCell>{getId(change.before.email)}</TableCell>
      <TableCell>{change.before.email}</TableCell>
      <TableCell>
        {getDiff(change.before.members, change.after.members)}
      </TableCell>
    </TableRow>
  </>
);

export default GroupUpdate;
