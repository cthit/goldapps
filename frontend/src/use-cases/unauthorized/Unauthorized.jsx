import { Typography } from "@mui/material";

const Unauthorized = () => {
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
      <Typography variant="p">
        You are not authorized to use this website
      </Typography>
      <Typography variant="p">
        For more information, please contact IT responsible at the IT student
        division or digIT
      </Typography>
    </div>
  );
};

export default Unauthorized;
