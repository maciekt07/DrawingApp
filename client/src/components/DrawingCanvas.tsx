import React, { useEffect, useRef, useState } from "react";
import { useDrawingStore } from "../store/store";
import { Drawing, Point } from "../types/types";

const DrawingCanvas: React.FC = () => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [currentPath, setCurrentPath] = useState<Point[]>([]);
  const [currentColor, setCurrentColor] = useState<string>("#000000");
  const { drawings, setDrawings, addDrawing } = useDrawingStore();
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    // Initialize WebSocket connection
    socketRef.current = new WebSocket("ws://localhost:8080/ws");

    socketRef.current.onmessage = (event) => {
      const newDrawing: Drawing = JSON.parse(event.data);
      addDrawing(newDrawing);
    };

    fetchDrawings();
    return () => {
      socketRef.current?.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    drawFetchedDrawings(drawings);
  }, [drawings]);

  const fetchDrawings = async () => {
    try {
      const response = await fetch("http://localhost:8080/drawings");
      const data: Drawing[] = await response.json();
      setDrawings(data);
    } catch (error) {
      console.error("Error fetching drawings:", error);
    }
  };

  const drawFetchedDrawings = (fetchedDrawings: Drawing[]) => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");

    if (ctx && canvas) {
      ctx.clearRect(0, 0, canvas.width, canvas.height); // clear canvas first
      fetchedDrawings.forEach(({ path, color }) => {
        if (path && Array.isArray(path) && path.length > 0) {
          ctx.strokeStyle = color;
          ctx.beginPath();
          ctx.moveTo(path[0].x, path[0].y);
          path.forEach((point) => {
            ctx.lineTo(point.x, point.y);
          });
          ctx.stroke();
        }
      });
    }
  };

  const startDrawing = (event: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");

    if (ctx && canvas) {
      setIsDrawing(true);
      const rect = canvas.getBoundingClientRect();
      const x = event.clientX - rect.left;
      const y = event.clientY - rect.top;

      ctx.beginPath();
      ctx.moveTo(x, y);
      setCurrentPath([{ x, y }]);
    }
  };

  const draw = (event: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isDrawing) return;

    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");

    if (ctx && canvas) {
      const rect = canvas.getBoundingClientRect();
      const x = event.clientX - rect.left;
      const y = event.clientY - rect.top;

      ctx.lineTo(x, y);
      ctx.stroke();

      setCurrentPath((prev) => [...prev, { x, y }]);
    }
  };

  const stopDrawing = () => {
    if (currentPath.length > 0) {
      const newDrawing = {
        path: currentPath,
        color: currentColor,
      };
      addDrawing(newDrawing);
      sendDrawingToWebSocket(newDrawing);
      saveDrawingToBackend(newDrawing);
    }
    setIsDrawing(false);
    setCurrentPath([]);
  };

  const saveDrawingToBackend = async (drawing: Drawing) => {
    try {
      await fetch("http://localhost:8080/drawings", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(drawing),
      });
    } catch (error) {
      console.error("Error saving drawing:", error);
    }
  };

  const sendDrawingToWebSocket = (drawing: Drawing) => {
    if (socketRef.current) {
      socketRef.current.send(JSON.stringify(drawing));
    }
  };

  return (
    <div>
      <div>
        <label>
          Color:
          <input
            type="color"
            value={currentColor}
            onChange={(e) => setCurrentColor(e.target.value)}
          />
        </label>
      </div>
      <canvas
        ref={canvasRef}
        width={800}
        height={600}
        style={{ border: "1px solid black" }}
        onMouseDown={startDrawing}
        onMouseMove={draw}
        onMouseUp={stopDrawing}
        onMouseLeave={stopDrawing}
      />
    </div>
  );
};

export default DrawingCanvas;
