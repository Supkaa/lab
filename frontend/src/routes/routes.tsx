import React from 'react';
import {createBrowserRouter, Route, Router, Routes} from "react-router-dom";
import App from '../App';
import LoginPage from "../pages/login.page";
import HomePage from "../pages/home.page";
import RegisterPage from "../pages/register.page";
import AdminPage from "../pages/admin.page";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <App/>,
        children: [
            {
                path: "/login",
                element: <LoginPage/>
            },
            {
                path: "/home",
                element: <HomePage/>
            },
            {
                path: "/register",
                element: <RegisterPage/>
            },
            {
                path: "/admin",
                element: <AdminPage/>
            }
        ]
    }
]);
