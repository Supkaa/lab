import React, {useState} from 'react';
import axios from "axios";

const LoginPage = () => {
    const [login, setLogin] = useState<string | null>(null);
    const [password, setPassword] = useState<string | null>(null);
    const [isError, setIsError] = useState<boolean>(false);

    const loginHandler = async (event: any) => {
        event.preventDefault();
        try {
            const response = await axios.post('http://localhost:3333/login',
                {
                email: login,
                password: password
            }, {
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            localStorage.setItem('user', JSON.stringify(response.data));
            window.location.href = ('/home');
        } catch (error) {
            setIsError(true)
        }
    }


    return (
        <section className="bg-gray-50 dark:bg-gray-900">
            <div className="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
                <a href="#" className="flex items-center mb-6 text-2xl font-semibold text-gray-900 dark:text-white">
                        Лабораторная 2
                </a>
                <div className="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
                    <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                        <h1 className="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
                            Войти в аккаунт
                        </h1>
                        <form className="space-y-4 md:space-y-6" action="#">
                            <div>
                                <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Email</label>
                                <input onChange={(event) => setLogin(event.target.value)} type="email" name="email" id="email" className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="test@gmail.com" required />
                            </div>
                            <div>
                                <label  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Пароль</label>
                                <input onChange={(event) => setPassword(event.target.value)} type="password" name="password" id="password" placeholder="••••••••" className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" required />
                            </div>
                            <div className="flex flex-col items-center justify-between">
                                <div className="flex items-start w-full">
                                    <button onClick={(event) => loginHandler(event)} className={'p-2 bg-purple-500 text-white rounded-md w-full'}>Войти</button>
                                </div>
                                {isError && (
                                    <p className={'text-red-500 mt-4'}>Неправильный логин или пароль</p>
                                )}
                            </div>
                            <p className="text-sm font-light text-gray-500 dark:text-gray-400 flex justify-between">
                               Нет аккаунта? <a href="/register" className="font-medium text-primary-600 hover:underline dark:text-primary-500">Регистрация</a>
                            </p>
                        </form>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default LoginPage;