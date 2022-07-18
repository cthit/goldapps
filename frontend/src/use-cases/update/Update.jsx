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
import {
  UserUpdate,
  UserAddition,
  UserDeletion,
  GroupUpdate,
  GroupAddition,
  GroupDeletion,
} from "./elements";

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
