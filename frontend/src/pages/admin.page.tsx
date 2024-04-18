import React from 'react';

const AdminPage = () => {
    const user = localStorage.getItem('user');
    let role = ''
    if (user) {
        role = JSON.parse(user).role;
    }

    return (
        <div>
            {role === 'admin' ? <h1>Админ страница</h1> :  <h1>Нет доступа</h1>}
        </div>
    );
};

export default AdminPage;