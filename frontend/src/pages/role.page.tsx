import React, { useEffect, useState } from "react";
import axios from "axios";
import { IRole } from "../interfaces/IRole";

const RolePage = () => {
  const [roleId, setRoleId] = useState<string[]>([]);
  const [endpointId, setEndpointId] = useState<string[]>([]);
  interface IEndpoint {
    id: string;
    url: string;
    method: string;
  }
  const [endpointUrl, setEndpointUrl] = useState<IEndpoint[]>([]);
  const [ids, setIds] = useState<string>();
  const [roles, setRoles] = useState<IRole[]>([]);
  let token: null | string = null;
  const userData = localStorage.getItem("user");
  let roleAdmin: string = "";
  if (userData) {
    token = JSON.parse(userData).token;
    if (!token) {
      roleAdmin = JSON.parse(userData).role;
      window.location.href = "/login";
    }
  }

  useEffect(() => {
    axios
      .get("http://localhost:3333/roles", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((response) => {
        setRoles(response.data);
      });
    axios
      .get("http://localhost:3333/endpoints", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((response) => {
        setEndpointUrl(response.data);
      });
  }, [token]);

  return (
    <div className={"flex flex-col items-center"}>
      <div className="relative overflow-x-auto mt-28 w-1/2">
        <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
          <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" className="px-6 py-3">
                Роль
              </th>
              <th scope="col" className="px-6 py-3">
                Действие
              </th>
            </tr>
          </thead>
          <tbody>
            {roles?.map((role) => (
              <tr
                key={role.name}
                className="bg-white border-b dark:bg-gray-800 dark:border-gray-700"
              >
                <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                  {role.name}
                </td>
                <td className="px-6 py-4 flex gap-4">
                  <button
                    onClick={() => {
                      axios
                        .delete(`http://localhost:3333/roles/${role.name}`)
                        .then(() => {
                          window.location.reload();
                        });
                    }}
                    className="bg-red-500 hover:bg-red-700 text-white font-bold p-1 rounded"
                  >
                    Удалить
                  </button>
                  <button
                    onClick={() => setRoleId([role.name])}
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold p-1 rounded"
                  >
                    Редактировать
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {roleId && (
          <div className={"flex gap-7 mt-4"}>
            <input
              type="text"
              onChange={(event) => setIds(event.target.value)}
              className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              placeholder="id эндпоинта"
            />
            <button
              onClick={() => {
                axios
                  .put(
                    `http://localhost:3333/roles/${roleId}`,
                    {
                      endpointIds: ids,
                    },
                    { headers: { Authorization: `Bearer ${token}` } },
                  )
                  .then(() => {
                    window.location.reload();
                  });
              }}
              className="bg-blue-500 rounded-s hover:bg-blue-700 text-white font-bold p-1 rounded"
            >
              Cохранить
            </button>
          </div>
        )}
      </div>
      <div className={"flex gap-7 mt-4 w-1/2"}>
        <div className={"flex flex-col w-full gap-4"}>
          <input
            type="text"
            onChange={(event) => setIds(event.target.value)}
            className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="Название роли"
          />
          <input
            type="text"
            onChange={(event) =>
              setEndpointId(event.target.value.split(",").map((e) => e.trim()))
            }
            className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="id эндпоинта"
          />
        </div>
        <button
          onClick={() => {
            axios
              .post(
                `http://localhost:3333/roles`,
                {
                  name: ids,
                  endpointIds: endpointId,
                },
                { headers: { Authorization: `Bearer ${token}` } },
              )
              .then(() => {
                window.location.reload();
              });
          }}
          className="bg-green-500 rounded-s hover:bg-green-700 text-white font-bold p-1 rounded"
        >
          Создать роль
        </button>
      </div>
      <table className="w-1/2 text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400 mt-4">
        <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
          <tr>
            <th className="px-6 py-3">Id</th>
            <th className="px-6 py-3">Метод</th>
          </tr>
        </thead>
        <tbody>
          {endpointUrl?.map((endpoint) => (
            <tr
              key={endpoint.id}
              className="bg-white border-b dark:bg-gray-800 dark:border-gray-700"
            >
              <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                {endpoint.id}
              </td>
              <td className="px-6 py-4 flex gap-4">{endpoint.url}</td>
              <td className="px-6 py-4 flex gap-4">{endpoint.method}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default RolePage;
