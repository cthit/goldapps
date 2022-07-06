import { BrowserRouter, Route, Routes } from "react-router-dom";
import { UpdateWizard } from "./use-cases";

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<UpdateWizard />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
