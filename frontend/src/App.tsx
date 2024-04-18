import React from 'react';
import './App.css';
import { Outlet } from 'react-router-dom';
import axios from "axios";
import Navbar from "./shared/navbar";

axios.defaults.baseURL = "http://localhost:3333";
axios.defaults.withCredentials = true;


function App() {
  return (
      <>
        <Navbar/>
        <Outlet />
      </>
  );
}

export default App;
