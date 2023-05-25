import { Typography } from "@mui/material";
import React from "react";

export default function Unauthorized() {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
      }}
    >
      <Typography variant="h2">Unauthorized!</Typography>
      <p>You are not authorized to use this website</p>
      <p>
        For more information, please contact IT responsible at the IT student
        division or digIT
      </p>
    </div>
  );
}
