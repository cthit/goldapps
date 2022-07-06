import {
  Box,
  Stepper,
  Step,
  StepLabel,
  Button,
  Card,
  CardContent,
  Typography,
} from "@mui/material";
import { DataGrid } from "@mui/x-data-grid";
import { useEffect, useState } from "react";
import Axios from "axios";

const steps = [
  { label: "Collect" },
  {
    label: "User Updates",
  },
  {
    label: "User Deletions",
  },
  {
    label: "User Additions",
  },
  {
    label: "Group Updates",
  },
  {
    label: "Group Deletions",
  },
  {
    label: "Group Additions",
  },
  { label: "Commit" },
];

const columns = [
  {
    field: "cid",
    headerName: "Cid",
    width: 100,
  },
  { field: "first_name", headerName: "First Name", width: 100 },
  { field: "second_name", headerName: "Last Name", width: 100 },
  { field: "mail", headerName: "E-mail", width: 200 },
];

const UpdateWizard = () => {
  const [activeStep, setActiveStep] = useState(0);
  const [rows, setRows] = useState({
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
      .then(res => setRows(res.data))
      .catch(err => console.log(err));
  }, []);

  return (
    <Card sx={{ margin: "1rem", padding: "1rem" }}>
      <CardContent>
        <Box sx={{ width: "100%" }}>
          <Stepper activeStep={activeStep} alternativeLabel>
            {steps.map(s => (
              <Step>
                <StepLabel>{s.label}</StepLabel>
              </Step>
            ))}
          </Stepper>
          <div
            style={{
              display: "flex",
              justifyContent: "center",
              marginTop: "1rem",
            }}
          >
            <Typography variant="h5">{steps[activeStep].label}</Typography>
          </div>
          <DataGrid
            sx={{ height: 570 }}
            rows={
              rows.userChanges.deletions
                ? rows.userChanges.deletions.map(d => ({ ...d, id: d.cid }))
                : []
            }
            columns={columns}
            pageSize={20}
            rowsPerPageOptions={[20, 50, 100]}
            checkboxSelection
          />
          <div style={{ display: "flex", "justify-content": "space-between" }}>
            <Button
              sx={{ marginTop: "2rem" }}
              onClick={() => setActiveStep(activeStep - 1)}
              variant="outlined"
              disabled={activeStep === 0}
            >
              Back
            </Button>
            <Button
              sx={{ marginTop: "2rem" }}
              onClick={() =>
                setActiveStep(Math.min(activeStep + 1, steps.length - 1))
              }
              variant="outlined"
            >
              {activeStep === steps.length - 1 ? "Commit" : "Next"}
            </Button>
          </div>
        </Box>
      </CardContent>
    </Card>
  );
};

export default UpdateWizard;
