const express = require("express");
const groups = require("./groups.json");
const users = require("./users.json");

app = express();

app.get("/api/goldapps/groups", (req, res) => {
  res.send(groups);
});

app.get("/api/goldapps/users", (req, res) => {
  res.send(users);
});

app.listen(8081);
