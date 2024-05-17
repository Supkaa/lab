import React from 'react';
import {createBrowserRouter, Route, Router, Routes} from "react-router-dom";
import App from '../App';
import LoginPage from "../pages/login.page";
import HomePage from "../pages/home.page";
import RegisterPage from "../pages/register.page";
import AdminPage from "../pages/admin.page";
import RolePage from "../pages/role.page";
import NewsPage from "../pages/news.page";
import NewsDetailPage from "../pages/news-detail.page";

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
                path: "/role",
                element: <RolePage/>
            },
            {
                path: "/register",
                element: <RegisterPage/>
            },
            {
                path: "/admin",
                element: <AdminPage/>
            },
            {
                path: "/news",
                element: <NewsPage/>
            },
            {
                path: "/news/:id",
                element: <NewsDetailPage/>
            }
        ]
    }
]);
