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
    set((state) => ({ drawings: [...state.drawings, drawing] })),
  setDrawings: (drawings) => set({ drawings }),
}));
