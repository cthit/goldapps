import { BrowserRouter, Route, Routes } from "react-router-dom";
import { UserProvider } from "./common/contexts/User.context";
import { Callback, Update, Unauthorized } from "./use-cases";
import history from "./utils/history";

const App = () => {
  return (
    <BrowserRouter history={history}>
      <UserProvider>
        <Routes>
          <Route path="/callback" element={<Callback />} />
          <Route path="/unauthorized" element={<Unauthorized />} />
          <Route path="/" element={<Update />} />
        </Routes>
      </UserProvider>
    </BrowserRouter>
  );
};

export default App;
