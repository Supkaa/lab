import React from 'react';
import './App.css';
import { Outlet } from 'react-router-dom';
import axios from "axios";

axios.defaults.baseURL = "http://localhost:3333";
axios.defaults.withCredentials = true;


function App() {
  return (
   <Outlet />
  );
}

export default App;
