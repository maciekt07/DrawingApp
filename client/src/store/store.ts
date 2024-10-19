import { create } from "zustand";
import { Drawing } from "../types/types";

interface DrawingStore {
  drawings: Drawing[];
  addDrawing: (drawing: Drawing) => void;
  setDrawings: (drawings: Drawing[]) => void;
}

export const useDrawingStore = create<DrawingStore>((set) => ({
  drawings: [],

  addDrawing: (drawing) =>
    set((state) => ({
      drawings: [...state.drawings, drawing],
    })),

  setDrawings: (drawings) => {
    if (Array.isArray(drawings)) {
      set({ drawings });
    } else {
      console.error("setDrawings: expected an array, received:", drawings);
    }
  },
}));
