import {
  Box,
  Button,
  ButtonGroup,
  Card,
  Checkbox,
  FormControl,
  FormControlLabel,
  IconButton,
  Modal,
  Paper,
  Radio,
  RadioGroup,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import MoreVertIcon from "@mui/icons-material/MoreVert";
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

export const formatEntry = entry => {
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

export const formatData = data => {
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

export const getAllIds = data => {
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

export const filterData = (data, selected) => {
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

const types = [
  { id: "user-update", label: "User Updates" },
  { id: "user-addition", label: "User Addition" },
  { id: "user-deletion", label: "User Deletion" },
  { id: "group-update", label: "Group Updates" },
  { id: "group-addition", label: "Group Addition" },
  { id: "group-deletion", label: "Group Deletion" },
];

const Update = () => {
  const [selected, setSelected] = useState([]);
  const [allSelected, setAllSelected] = useState(true);
  const [numEntries, setNumEntries] = useState(0);
  const [provider, setProvider] = useState("gamma");
  const [providerFile, setProviderFile] = useState("");
  const [consumer, setConsumer] = useState("json");
  const [consumerFile, setConsumerFile] = useState("data.json");
  const [typeSelectOpen, setTypeSelectOpen] = useState(false);
  const [typesSelected, setTypesSelected] = useState(types.map(t => t.id));

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

  const handleTypeSelect = id => {
    if (typesSelected.includes(id)) {
      setTypesSelected(typesSelected.filter(e => e !== id));
    } else {
      setTypesSelected([...typesSelected, id]);
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
      <Box sx={{ paddingTop: "2rem", paddingLeft: "2rem" }}>
        <ButtonGroup>
          <Typography sx={{ paddingTop: "1rem", paddingRight: "1rem" }}>
            From:
          </Typography>
          <RadioGroup
            sx={{ display: "flex", flexDirection: "row" }}
            defaultValue="gamma"
            onChange={e => setProvider(e.target.value)}
            value={provider}
            aria-labelledby="provider-radio-button"
            name="provider-radio-button"
          >
            <FormControlLabel
              value="gapps"
              control={<Radio />}
              label="G-suit"
            />
            <FormControlLabel value="gamma" control={<Radio />} label="Gamma" />
            <FormControlLabel value="ldap" control={<Radio />} label="LDAP" />
            <FormControlLabel value="json" control={<Radio />} label=".json" />
            <TextField
              variant="standard"
              value={providerFile}
              onChange={e => setProviderFile(e.target.value)}
              disabled={provider !== "json"}
            />
          </RadioGroup>
        </ButtonGroup>
        <ButtonGroup sx={{ paddingTop: "1rem" }}>
          <Typography sx={{ paddingTop: "1rem", paddingRight: "2.2rem" }}>
            To:
          </Typography>
          <RadioGroup
            sx={{ display: "flex", flexDirection: "row" }}
            defaultValue="json"
            onChange={e => setConsumer(e.target.value)}
            value={consumer}
            aria-labelledby="consumer-radio-button"
            name="consumer-radio-button"
          >
            <FormControlLabel
              value="gapps"
              control={<Radio />}
              label="G-suit"
            />

            <FormControlLabel value="json" control={<Radio />} label=".json" />
            <TextField
              variant="standard"
              value={consumerFile}
              onChange={e => setConsumerFile(e.target.value)}
              disabled={consumer !== "json"}
            />
          </RadioGroup>
        </ButtonGroup>
        <div>
          <Button
            sx={{ marginTop: "1rem", textTransform: "none" }}
            variant="contained"
          >
            Collect suggestions
          </Button>
        </div>
      </Box>
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
              <TableCell>
                Type
                <IconButton
                  onClick={() => {
                    setTypeSelectOpen(true);
                  }}
                >
                  <MoreVertIcon />
                </IconButton>
                <Modal
                  open={typeSelectOpen}
                  onClose={() => {
                    console.log("Closing modal!");
                    setTypeSelectOpen(false);
                  }}
                  aria-labelledby="parent-modal-title"
                  aria-describedby="parent-modal-description"
                  style={{
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                  }}
                >
                  <Card sx={{ padding: "2rem" }}>
                    <FormControl value={typesSelected}>
                      {types.map(t => (
                        <FormControlLabel
                          control={
                            <Checkbox
                              onChange={() => handleTypeSelect(t.id)}
                              checked={typesSelected.includes(t.id)}
                              value={t.id}
                            />
                          }
                          label={t.label}
                        />
                      ))}
                    </FormControl>
                  </Card>
                </Modal>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {typesSelected.includes("user-changes") &&
            data.userChanges.userUpdates
              ? data.userChanges.userUpdates.map(e => (
                  <UserUpdate
                    key={e.id}
                    change={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {typesSelected.includes("user-addition") &&
            data.userChanges.additions
              ? data.userChanges.additions.map(e => (
                  <UserAddition
                    key={e.id}
                    addition={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {typesSelected.includes("user-deletion") &&
            data.userChanges.deletions
              ? data.userChanges.deletions.map(e => (
                  <UserDeletion
                    key={e.id}
                    deletion={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {typesSelected.includes("group-update") &&
            data.groupChanges.groupUpdates
              ? data.groupChanges.groupUpdates.map(e => (
                  <GroupUpdate
                    key={e.id}
                    change={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {typesSelected.includes("group-addition") &&
            data.groupChanges.additions
              ? data.groupChanges.additions.map(e => (
                  <GroupAddition
                    key={e.id}
                    addition={e}
                    selected={selected}
                    onChange={onCheck}
                  />
                ))
              : null}
            {typesSelected.includes("group-deletion") &&
            data.groupChanges.deletions
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
