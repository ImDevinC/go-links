import React from 'react';
import './App.css';
import { CssBaseline, CssVarsProvider } from '@mui/joy';
import { Home } from './views/home';

function App() {
  return (
    <CssVarsProvider>
      <CssBaseline />
      <Home />
    </CssVarsProvider>
  );
}

export default App;
