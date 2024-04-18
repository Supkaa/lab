import React, {useEffect, useState} from 'react';
import {IUser} from "../interfaces/IUser";
import axios from "axios";
import {useNavigate} from "react-router-dom";

const HomePage = () => {
    const [users, setUsers] = useState<IUser[] | null>(null);
    const navigate = useNavigate();
    const [user, _] = useState<IUser>({
        name: '',
        email: '',
        password: '',
        role: ''
    });
    const [isOpenEditField, setIsOpenEditField] = useState<null | string>(null);
    let token:null | string = null;
    let role:string = '';
    const userData = localStorage.getItem('user');
    if (userData) {
        role = JSON.parse(userData).role;
        token = JSON.parse(userData).token;
    }

    const handleEdit = (email: string) => {
        axios.put(`http://localhost:3333/users/${email}`, {
            name: user.name,
            password: user.password,
            email: user.email,
            role: user.role
        }, {headers: {Authorization: `Bearer ${token}`}}).then(() => {
            setIsOpenEditField(null);
        })
    }


    const getUsers = async () => {
        try {
            if (!token) {
                navigate('/login');
            }
            const response = await axios.get('http://localhost:3333/users', {
                headers: {
                        Authorization: `Bearer ${token}`
                }
            });
            setUsers(response.data);
        } catch (error) {
            console.error('Ошибка при получении пользователей:', error);
        }
    };
    useEffect(() => {
       getUsers();
    }, [isOpenEditField])

    return (
        <div className={'flex flex-col w-full items-center mt-28'}>
            <div className="relative overflow-x-auto shadow-md sm:rounded-lg w-1/2">
                <h1 className="text-2xl font-bold mb-4">Пользователи</h1>
                <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                    <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    <tr>
                        <th scope="col" className="px-6 py-3">
                            Еmail
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Имя
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Роль
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Действие
                        </th>
                    </tr>
                    </thead>
                    <tbody>
                    {users?.map((user) => (
                        <tr key={user.email} className="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                            <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                                {user.name}
                            </td>
                            <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                                {user.email}
                            </td>
                            <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                                {user.role}
                            </td>
                            {role === 'admin' && (
                                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium gap-6 flex">
                                    <button className="text-purple-400" onClick={() => setIsOpenEditField(user.email)}>Редактировать</button>
                                </td>
                            )}
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
                {isOpenEditField && (
                    <div className={'w-1/2 flex flex-col items-center gap-3 mt-4'}>
                        <input
                            type="text"
                            name="name"
                            onChange={(event) => user.name = (event.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                            placeholder="Имя"
                        />
                        <input
                            type="email"
                            name="email"
                            onChange={(event) => user.email = (event.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                            placeholder="Email"
                        />
                        <input
                            type="password"
                            name="password"
                            onChange={(event) => user.password = (event.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                            placeholder="Пароль"
                        />
                        <select
                            name="role"
                            onChange={(event) => user.role = (event.target.value)}
                            className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                        >
                            <option value="user">Пользователь</option>
                            <option value="admin">Администратор</option>
                        </select>
                        <button type="submit" onClick={() => handleEdit(isOpenEditField)}>Сохранить</button>
                    </div>
                )}
        </div>
    );
};

export default HomePage;