import DrawingCanvas from "./components/DrawingCanvas";

const App: React.FC = () => {
  return (
    <div>
      <div style={{ textAlign: "center" }}>
        <h1>Real-Time Drawing App</h1>
        <h2>
          Open it in multiple browser windows to draw together in real-time.
        </h2>
        <DrawingCanvas />
      </div>
    </div>
  );
};

export default App;
