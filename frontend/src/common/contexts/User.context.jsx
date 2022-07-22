import axios from "axios";
import { useEffect, useState } from "react";

const { createContext } = require("react");

export const UserContext = createContext({ admin: false });

export const UserProvider = ({ children }) => {
  const [user, setUser] = useState({ admin: false });

  useEffect(() => {
    if (
      window.location.pathname === "/callback" ||
      window.location.pathname === "/unauthorized"
    ) {
      return;
    }
    axios
      .get("/api/checkLogin")
      .then(res => setUser({ admin: true }))
      .catch(err => {
        console.log({ ...err });
        if (err.response.status === 401) {
          window.location.href = err.response.data;
        }
      });
  }, []);

  return (
    <UserContext.Provider value={[user, setUser]}>
      {children}
    </UserContext.Provider>
  );
};
