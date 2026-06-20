# 📷 Foto Fácil

A node-based image processing platform that executes batch image transformations using a Directed Acyclic Graph (DAG) system. Built with Electron, React, TypeScript, and Go.

### 💡 About the Project
Foto Fácil was developed to automate repetitive image editing tasks. Instead of manual processing one by one, the user can visually build a pipeline of processing blocks (nodes) to perform actions like cropping, rotating, brightness adjustment, color conversions, and AI-driven background removal, applying the entire sequence to multiple images or entire folders automatically.

### ✨ Key Features
The application provides an intuitive flow-canvas where you can drag, connect, configure, and execute image processing pipelines.

#### 📦 Node Library
* **Image Input**: Select and load one or more images from disk, or target an entire directory to batch-process all contained image files automatically.
* **Image Output**: Select the output directory where all processed images will be saved. Supports PNG format to preserve alpha channel transparency.
* **Brightness & Contrast**: Adjust the exposure level of the images dynamically on a scale of -255 to 255.
* **Color Space**: Convert standard color images to Grayscale using the luminance formula (0.299R + 0.587G + 0.114B).
* **Crop/Resize**: Crop a sub-region (X, Y, Width, Height) of the images and/or resize them to target dimensions.
* **Rotate/Flip**: Rotate images in 90° increments (90°, 180°, 270°) and flip them horizontally, vertically, or both.
* **IA/Machine Learning**: Advanced smart tools including background removal and intelligent contrast boost, featuring interactive tolerance control.
* **Comparação Visual**: A visual comparison slider that allows dragging to see the "before" and "after" state of the image side-by-side.

#### ⚙️ Workspace & System Features
* **Excel-style Multi-Flows**: Create, rename (double-click), and delete multiple flow worksheets at the bottom bar.
* **SQLite Persistence**: All flows, including nodes, coordinates, settings, and connections, are automatically saved and restored from a local SQLite database (`foto-facil.db`).
* **Dynamic Dark/Light Mode**: Smooth transitions between a sleek developer dark mode and a high-contrast light mode with custom node pastel colors.
* **Context Menus**: Right-click nodes or edges to delete them from the workspace.

### 🚀 How to Run

#### Prerequisites
* [Go](https://go.dev/doc/install) (Go 1.18+ recommended)
* [Node.js](https://nodejs.org/) (v16+ recommended)

#### 1. Start the Go Backend
```bash
cd backend
go run cmd/main.go
```
The server will start listening on port `:8080` and expose the WebSocket connection.

#### 2. Start the Frontend & Electron App
```bash
cd frontend
npm install
npm run dev
```
This boots up the Vite development server and launches the Electron application automatically.

Developed by @nathanhgo.