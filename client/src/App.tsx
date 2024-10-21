import DrawingCanvas from "./components/DrawingCanvas";

const App: React.FC = () => {
  return (
    <div>
      <div style={{ textAlign: "center" }}>
        <h1>Real-Time Drawing App</h1>
        <h2>Open this app in multiple tabs to draw simultaneously.</h2>
        <DrawingCanvas />
      </div>
    </div>
  );
};

export default App;
