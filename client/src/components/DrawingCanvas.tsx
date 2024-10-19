import React, { useEffect, useRef } from "react";
import { useDrawingStore } from "../store/store";
import { Drawing } from "../types/types";

// TODO: rewrite this completely
const DrawingCanvas: React.FC = () => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const { drawings, addDrawing, setDrawings } = useDrawingStore();
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const fetchDrawings = async () => {
      try {
        const response = await fetch("http://localhost:8080/drawings");
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const fetchedDrawings: Drawing[] = await response.json();
        setDrawings(fetchedDrawings);
      } catch (error) {
        console.error("Failed to fetch drawings:", error);
      }
    };

    fetchDrawings();
  }, [setDrawings]);

  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");

    if (ctx && canvas) {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      drawings.forEach(({ x, y, color }) => {
        ctx.fillStyle = color;
        ctx.fillRect(x, y, 5, 5);
      });
    }
  }, [drawings]);

  // WebSocket connection
  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socketRef.current = socket;

    socket.onopen = () => {
      console.log("WebSocket connection established");
    };

    socket.onmessage = (event) => {
      const message: Drawing = JSON.parse(event.data);
      addDrawing(message);
      saveDrawingToDB(message); // Save drawing to DB on receiving a message
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    socket.onclose = () => {
      console.error("WebSocket connection closed");
    };

    return () => {
      socket.close();
    };
  }, [addDrawing]);

  const handleMouseMove = (event: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const color = "black"; // TODO: make this dynamic

    const drawing: Drawing = { x, y, color };
    addDrawing(drawing);
    drawOnCanvas(x, y, color);

    // Send to WebSocket only if it is open
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify(drawing));
    } else {
      console.warn(
        "WebSocket is not open. Current state: " + socketRef.current?.readyState
      );
    }
  };

  const drawOnCanvas = (x: number, y: number, color: string) => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");
    if (ctx) {
      ctx.fillStyle = color;
      ctx.fillRect(x, y, 5, 5);
    }
  };

  // save drawing to the database
  const saveDrawingToDB = async (drawing: Drawing) => {
    await fetch("http://localhost:8080/drawings", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(drawing),
    });
  };

  return (
    <canvas
      ref={canvasRef}
      width={800}
      height={600}
      style={{ border: "1px solid black" }}
      onMouseMove={handleMouseMove}
    />
  );
};

export default DrawingCanvas;
