import axios from "axios";
import { useContext, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { UserContext } from "../../common/contexts/User.context";

const Callback = () => {
  const [, setUser] = useContext(UserContext);
  const navigate = useNavigate();
  useEffect(() => {
    const params = new URL(window.location.href).searchParams;
    axios
      .get("/api/authenticate", {
        params: {
          code: params.get("code"),
        },
      })
      .then(() => {
        setUser({ admin: true });
        navigate("/");
      })
      .catch(() => {
        setUser({ admin: false });
        navigate("/unauthorized");
      });
  }, [setUser, navigate]);
  return <div></div>;
};

export default Callback;
