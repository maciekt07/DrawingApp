import React, { useEffect, useRef, useState } from "react";
import { useDrawingStore } from "../store/store";
import { Drawing, Point } from "../types/types";

const DrawingCanvas: React.FC = () => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [currentPath, setCurrentPath] = useState<Point[]>([]);
  const [currentColor, setCurrentColor] = useState<string>("#ffffff");
  const { drawings, setDrawings, addDrawing, clearDrawings } =
    useDrawingStore();
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const initSocket = () => {
      socketRef.current = new WebSocket("ws://localhost:8080/ws");
      socketRef.current.onmessage = (event) => {
        const newDrawing = JSON.parse(event.data);
        if (newDrawing.type === "clear") {
          clearDrawings(); // clear drawings when event is received
        } else {
          addDrawing(newDrawing);
        }
      };
    };

    const fetchDrawings = async () => {
      try {
        const response = await fetch("/api/drawings");
        const data: Drawing[] = await response.json();
        setDrawings(data);
      } catch (error) {
        console.error("Error fetching drawings:", error);
      }
    };

    initSocket();
    fetchDrawings();

    return () => socketRef.current?.close();
  }, [addDrawing, clearDrawings, setDrawings]);

  useEffect(() => {
    const ctx = canvasRef.current?.getContext("2d");
    if (ctx && canvasRef.current) {
      ctx.clearRect(0, 0, canvasRef.current.width, canvasRef.current.height);
      drawings.forEach(({ path, color }) => {
        if (path.length) {
          ctx.strokeStyle = color;
          ctx.beginPath();
          ctx.moveTo(path[0].x, path[0].y);
          path.forEach(({ x, y }) => ctx.lineTo(x, y));
          ctx.stroke();
        }
      });
    }
  }, [drawings]);

  const startDrawing = ({
    clientX,
    clientY,
  }: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");
    if (ctx && canvas) {
      setIsDrawing(true);
      const { left, top } = canvas.getBoundingClientRect();
      ctx.strokeStyle = currentColor;
      ctx.beginPath();
      ctx.moveTo(clientX - left, clientY - top);
      setCurrentPath([{ x: clientX - left, y: clientY - top }]);
    }
  };

  const draw = ({ clientX, clientY }: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isDrawing) return;
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");
    if (ctx && canvas) {
      const { left, top } = canvas.getBoundingClientRect();
      ctx.lineTo(clientX - left, clientY - top);
      ctx.stroke();
      setCurrentPath((prev) => [
        ...prev,
        { x: clientX - left, y: clientY - top },
      ]);
    }
  };

  const stopDrawing = () => {
    if (currentPath.length) {
      const newDrawing = { path: currentPath, color: currentColor };
      addDrawing(newDrawing);
      socketRef.current?.send(JSON.stringify(newDrawing));
      setCurrentPath([]);
    }
    setIsDrawing(false);
  };

  const clearCanvas = () => {
    clearDrawings();
    socketRef.current?.send(JSON.stringify({ type: "clear" })); // Send clear event to all clients
  };

  return (
    <div>
      <div className="color-picker-container">
        <label className="color-label">
          Color:
          <input
            type="color"
            value={currentColor}
            onChange={(e) => setCurrentColor(e.target.value)}
            className="color-input"
          />
        </label>
      </div>
      <button onClick={clearCanvas}>Clear</button>
      <div className="canvas-container">
        <canvas
          ref={canvasRef}
          width={800}
          height={600}
          onMouseDown={startDrawing}
          onMouseMove={draw}
          onMouseUp={stopDrawing}
          onMouseLeave={stopDrawing}
        />
      </div>
    </div>
  );
};

export default DrawingCanvas;
