import {
  Box,
  Button,
  Checkbox,
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
import { getId } from "../../utils/utils";

const formatEntry = entry => {
  let keys = Object.keys(entry);
  if (keys.includes("cid")) {
    return { ...entry, id: entry.cid };
  }
  if (keys.includes("members")) {
    return { ...entry, id: getId(entry.email) };
  }
  keys = Object.keys(entry.before);
  if (keys.includes("cid")) {
    return { ...entry, id: entry.before.cid };
  }
  if (keys.includes("members")) {
    return { ...entry, id: getId(entry.before.email) };
  }
  return entry;
};

const formatData = data => {
  const ids = [];
  for (const change in data) {
    for (const type in data[change]) {
      for (const i in data[change][type]) {
        data[change][type][i] = formatEntry(data[change][type][i]);
        ids.push(data[change][type][i].id);
      }
    }
  }
  return [data, ids];
};

const getAllIds = data => {
  const ids = [];
  for (const change in data) {
    for (const type in data[change]) {
      for (const i in data[change][type]) {
        ids.push(data[change][type][i].id);
      }
    }
  }
  return ids;
};

const filterData = (data, selected) => {
  for (const change in data) {
    for (const type in data[change]) {
      if (data[change][type] === null) {
        continue;
      }
      data[change][type] = data[change][type].filter(e =>
        selected.includes(e.id),
      );
      if (data[change][type].length === 0) {
        delete data[change][type];
      }
    }
  }
  return data;
};

const Update = () => {
  const [selected, setSelected] = useState([]);
  const [allSelected, setAllSelected] = useState(true);
  const [numEntries, setNumEntries] = useState(0);
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
      .then(res => {
        const [data, ids] = formatData(res.data);
        setData(data);
        setSelected(ids);
        setAllSelected(true);
        setNumEntries(ids.length);
      })
      .catch(err => console.log(err));
  }, []);

  const onCheckAll = () => {
    if (allSelected) {
      setAllSelected(false);
      setSelected([]);
    } else {
      setAllSelected(true);
      setSelected(getAllIds(data));
    }
  };

  const onCheck = id => {
    if (selected.includes(id)) {
      setSelected(selected.filter(e => e !== id));
      setAllSelected(false);
    } else {
      setSelected([...selected, id]);
      setAllSelected(selected.length === numEntries - 1);
    }
  };

  const onCommit = () => {
    const selectedData = filterData(JSON.parse(JSON.stringify(data)), selected);
    const allIds = getAllIds(data);
    const unselectedData = filterData(
      JSON.parse(JSON.stringify(data)),
      allIds.filter(id => !selected.includes(id)),
    );

    Axios.post("/api/commit", selectedData)
      .then(res => {
        setData(unselectedData);
        setSelected([]);
        setAllSelected(false);
      })
      .catch(err => {
        console.log(err);
      });
  };

  return (
    <Box sx={{ width: "100%" }}>
      <Box sx={{ padding: "2rem" }}>
        <TableContainer component={Paper}>
          <TableHead>
            <TableRow>
              <TableCell padding="checkbox">
                <Checkbox checked={allSelected} onChange={onCheckAll} />
              </TableCell>
              <TableCell>Id</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>E-mail(s)</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.userChanges.userUpdates
              ? data.userChanges.userUpdates.map(e => (
                  <UserUpdate
                    key={e.id}
                    change={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {data.userChanges.additions
              ? data.userChanges.additions.map(e => (
                  <UserAddition
                    key={e.id}
                    addition={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {data.userChanges.deletions
              ? data.userChanges.deletions.map(e => (
                  <UserDeletion
                    key={e.id}
                    deletion={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {data.groupChanges.groupUpdates
              ? data.groupChanges.groupUpdates.map(e => (
                  <GroupUpdate
                    key={e.id}
                    change={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {data.groupChanges.additions
              ? data.groupChanges.additions.map(e => (
                  <GroupAddition
                    key={e.id}
                    addition={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {data.groupChanges.deletions
              ? data.groupChanges.deletions.map(e => (
                  <GroupDeletion
                    key={e.id}
                    deletion={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
          </TableBody>
        </TableContainer>
        <div style={{ display: "flex", justifyContent: "flex-end" }}>
          <Button
            sx={{ marginTop: "2rem" }}
            variant="outlined"
            onClick={onCommit}
          >
            Commit
          </Button>
        </div>
      </Box>
    </Box>
  );
};

export default Update;
