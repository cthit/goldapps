import {
  Box,
  Button,
  Paper,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import { useEffect, useState } from "react";
import Axios from "axios";
import "./Update.css";

//id [cid], name [first_name 'nick' second_name], email [email]
const UserUpdate = ({ change }) => (
  <>
    <TableRow>
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

//id [cid], name [first_name 'nick' second_name], email [email]
const UserAddition = ({ addition }) => (
  <>
    <TableRow>
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

const getId = email => email.substr(0, email.search("@"));
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
const GroupUpdate = ({ change }) => (
  <>
    <TableRow>
      <TableCell>{getId(change.before.email)}</TableCell>
      <TableCell>{change.before.email}</TableCell>
      <TableCell>
        {getDiff(change.before.members, change.after.members)}
      </TableCell>
    </TableRow>
  </>
);

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

const GroupDeletion = ({ deletion }) => (
  <>
    <TableRow>
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

const Update = () => {
  const [data, setData] = useState({
    userChanges: {
      userDeletions: null,
      additions: null,
      deletions: null,
    },
    groupChanges: {
      groupUpdates: null,
      additions: null,
      deletions: null,
    },
  });

  useEffect(() => {
    Axios.get("/api/suggestions")
      .then(res => setData(res.data))
      .catch(err => console.log(err));
  }, []);

  return (
    <Box sx={{ width: "100%" }}>
      <Box sx={{ padding: "2rem" }}>
        <TableContainer component={Paper}>
          <TableHead>
            <TableRow>
              <TableCell>Id</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>E-mail(s)</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.userChanges.userUpdates
              ? data.userChanges.userUpdates.map(e => (
                  <UserUpdate key={e.before.cid} change={e} />
                ))
              : null}
            {data.userChanges.additions
              ? data.userChanges.additions.map(e => (
                  <UserAddition addition={e} />
                ))
              : null}
            {data.userChanges.deletions
              ? data.userChanges.deletions.map(e => (
                  <UserDeletion deletion={e} />
                ))
              : null}
            {data.groupChanges.groupUpdates
              ? data.groupChanges.groupUpdates.map(e => (
                  <GroupUpdate change={e} />
                ))
              : null}
            {data.groupChanges.additions
              ? data.groupChanges.additions.map(e => (
                  <GroupAddition addition={e} />
                ))
              : null}
            {data.groupChanges.deletions
              ? data.groupChanges.deletions.map(e => (
                  <GroupDeletion deletion={e} />
                ))
              : null}
          </TableBody>
        </TableContainer>
        <div style={{ display: "flex", justifyContent: "flex-end" }}>
          <Button sx={{ marginTop: "2rem" }} variant="outlined">
            Commit
          </Button>
        </div>
      </Box>
    </Box>
  );
};

export default Update;
