// preload.ts — expõe APIs seguras do Electron para o React
import { contextBridge, ipcRenderer } from 'electron';

contextBridge.exposeInMainWorld('electronAPI', {
  openFileDialog: () => ipcRenderer.invoke('dialog:openFile'),
  openDirectoryDialog: () => ipcRenderer.invoke('dialog:openDirectory'),
});

window.addEventListener('DOMContentLoaded', () => {
  console.log('Foto Fácil Electron Preload carregado (com IPC)');
});
