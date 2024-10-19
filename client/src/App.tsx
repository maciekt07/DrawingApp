// src/App.tsx
import React from "react";
import DrawingCanvas from "./components/DrawingCanvas";

const App: React.FC = () => {
  return (
    <div>
      <h1>Real-Time Drawing App</h1>
      <DrawingCanvas />
    </div>
  );
};

export default App;
