import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { INews } from "../interfaces/INews";
import { Link } from 'react-router-dom';

const NewsPage = () => {
    const [news, setNews] = useState<INews[]>([]);

    useEffect(() => {
        axios.get("http://localhost:3333/news")
            .then(res => {
                const fetchedNews = res.data;
                const updatedNews = fetchedNews.map((item: INews) => {
                    const storedLikes = localStorage.getItem(`likes-${item.id}`);
                    const storedDislikes = localStorage.getItem(`dislikes-${item.id}`);
                    if (storedLikes) {
                        item.likes = parseInt(storedLikes, 10);
                    }
                    if (storedDislikes) {
                        item.dislikes = parseInt(storedDislikes, 10);
                    }
                    return item;
                });
                setNews(updatedNews);
            });
    }, []);

    const handleLike = (id: string) => {
        const updatedNews = news.map((item) => {
            if (item.id === id) {
                const newLikes = item.likes + 1;
                localStorage.setItem(`likes-${id}`, newLikes.toString());
                return { ...item, likes: newLikes };
            }
            return item;
        });
        setNews(updatedNews);
    };

    const handleDislike = (id: string) => {
        const updatedNews = news.map((item) => {
            if (item.id === id) {
                const newDislikes = item.dislikes + 1;
                localStorage.setItem(`dislikes-${id}`, newDislikes.toString());
                return { ...item, dislikes: newDislikes };
            }
            return item;
        });
        setNews(updatedNews);
    };

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-4xl font-bold mb-6 text-center">Последние новости</h1>
            <div id="news-container"
                 className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
                {news.map((n) => (
                    <div key={n.id} className="bg-white rounded-lg shadow-md overflow-hidden p-4">
                        <Link to={`/news/${n.id}`}>
                            <h2 className="text-xl font-bold mb-2">{n.title}</h2>
                            {n.image ? (
                                <img style={{ height: 200, backgroundColor: 'gray', width: '100%' }} src={n.image} alt={n.title} />
                            ) : (
                                <div style={{ height: 200, backgroundColor: 'gray', width: '100%' }}></div>
                            )}
                        </Link>
                        <p className="text-gray-600">{n.summary}</p>
                        <p className="text-gray-600">{new Date(n.createdAt).toLocaleDateString()}</p>
                        <p className="text-gray-600 flex gap-2">{localStorage.getItem(`views-${n.id}`) || 0}
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="w-6 h-6">
                                <path strokeLinecap="round" strokeLinejoin="round"
                                      d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"/>
                                <path strokeLinecap="round" strokeLinejoin="round"
                                      d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"/>
                            </svg>
                        </p>
                        <p className="text-gray-600 flex gap-2">
                            {n.likes}
                            <svg
                                onClick={() => handleLike(n.id)}
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                strokeWidth="1.5"
                                stroke="currentColor"
                                className="w-6 h-6 cursor-pointer"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round"
                                      d="M6.633 10.25c.806 0 1.533-.446 2.031-1.08a9.041 9.041 0 0 1 2.861-2.4c.723-.384 1.35-.956 1.653-1.715a4.498 4.498 0 0 0 .322-1.672V2.75a.75.75 0 0 1 .75-.75 2.25 2.25 0 0 1 2.25 2.25c0 1.152-.26 2.243-.723 3.218-.266.558.107 1.282.725 1.282m0 0h3.126c1.026 0 1.945.694 2.054 1.715.045.422.068.85.068 1.285a11.95 11.95 0 0 1-2.649 7.521c-.388.482-.987.729-1.605.729H13.48c-.483 0-.964-.078-1.423-.23l-3.114-1.04a4.501 4.501 0 0 0-1.423-.23H5.904m10.598-9.75H14.25M5.904 18.5c.083.205.173.405.27.602.197.4-.078.898-.523.898h-.908c-.889 0-1.713-.518-1.972-1.368a12 12 0 0 1-.521-3.507c0-1.553.295-3.036.831-4.398C3.387 9.953 4.167 9.5 5 9.5h1.053c.472 0 .745.556.5.96a8.958 8.958 0 0 0-1.302 4.665c0 1.194.232 2.333.654 3.375Z"/>
                            </svg>
                        </p>
                        <p className="text-gray-600 flex gap-2">
                            {n.dislikes}
                            <svg
                                onClick={() => handleDislike(n.id)}
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                strokeWidth="1.5"
                                stroke="currentColor"
                                className="w-6 h-6 cursor-pointer"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round"
                                      d="M7.498 15.25H4.372c-1.026 0-1.945-.694-2.054-1.715a12.137 12.137 0 0 1-.068-1.285c0-2.848.992-5.464 2.649-7.521C5.287 4.247 5.886 4 6.504 4h4.016a4.5 4.5 0 0 1 1.423.23l3.114 1.04a4.5 4.5 0 0 0 1.423.23h1.294M7.498 15.25c.618 0 .991.724.725 1.282A7.471 7.471 0 0 0 7.5 19.75 2.25 2.25 0 0 0 9.75 22a.75.75 0 0 0 .75-.75v-.633c0-.573.11-1.14.322-1.672.304-.76.93-1.33 1.653-1.715a9.04 9.04 0 0 0 2.86-2.4c.498-.634 1.226-1.08 2.032-1.08h.384m-10.253 1.5H9.7m8.075-9.75c.01.05.027.1.05.148.593 1.2.925 2.55.925 3.977 0 1.487-.36 2.89-.999 4.125m"/>
                            </svg>
                        </p>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default NewsPage;