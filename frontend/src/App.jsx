import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Update } from "./use-cases";

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Update />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
