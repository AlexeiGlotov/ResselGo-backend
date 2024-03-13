
import {axiosInstanceWithJWT} from '../api/axios';
import React, { useState, useEffect } from 'react';
import {toast} from "react-toastify";


function LicenseKeys() {

    const [cheats, setCheats] = useState([]); // Данные для cheats
    const [users, setUsers] = useState([]); // Данные для user logins and roles
    const [loading, setLoading] = useState(true);
    const [subsCount, setSubsCount] = useState(1);
    const [daysCount, setDaysCount] = useState(25);
    const [selectedCheat, setSelectedCheat] = useState('');
    const [selectedUser, setSelectedUser] = useState('');
    const [keys, setKeys] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [licenseKeys, setLicenseKeys] = useState([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [filter, setFilter] = useState('');

    const handleCheatChange = (event) => {
        setSelectedCheat(event.target.value);
    };

    const handleUserChange = (event) => {
        setSelectedUser(event.target.value);
    };


    const handleSubsChange = (event) => {
        setSubsCount(event.target.value);
    };

    const handleDaysChange = (event) => {
        setDaysCount(event.target.value);
    };

    const handleSubmit = async () => {

        setIsLoading(true);
        const dataToSend = {
            count_keys: parseInt(subsCount),
            ttl_cheat: parseInt(daysCount),
            cheat:selectedCheat,
            holder:selectedUser
        };


        try {
            const [licenseResponse] = await Promise.all([
                axiosInstanceWithJWT.post('/api/license-keys/',dataToSend),
            ]);
            setKeys(licenseResponse.data.keys)
        } catch (error) {
            toast.error(`error: ${error.message}`);
        } finally {
            setIsLoading(false);
        }

    };

    useEffect(() => {
        const fetchLicenseKeys = async () => {
            try {
                const response = await axiosInstanceWithJWT.get(`api/license-keys/?page=${currentPage}`)
                setLicenseKeys(response.data.keys);
            } catch (error) {
                toast.error(`error: ${error.message}`);
            } finally {
                setIsLoading(false);
            }
        };

        // Вызов функции
        fetchLicenseKeys();

    }, [currentPage]);

    useEffect(() => {
        const fetchDatas = async () => {
            try {
                const [cheatsResponse, usersResponse] = await Promise.all([
                    axiosInstanceWithJWT.get('/api/cheats/'),
                    axiosInstanceWithJWT.get('/api/users/getUsers')
                ]);
                setCheats(cheatsResponse.data.cheats);
                setUsers(usersResponse.data.users); // Убедитесь, что здесь правильно обрабатывается ответ
            } catch (error) {
                toast.error(`error: ${error.message}`);
            } finally {
                setLoading(false);
            }
        };


        fetchDatas();
    }, []);



    if (loading) {
        return <div>Loading...</div>;
    }

    const renderGenerateCheat = () => {
        return (
        <div>
        <select value={selectedCheat} onChange={handleCheatChange}>
            <option value="">select cheat</option>
            {cheats.map((cheat) => (
                <option key={cheat.id} value={cheat.name}>
                    {cheat.name}
                </option>
            ))}
        </select>


        <select value={selectedUser} onChange={handleUserChange}>
            <option value="">select user</option>
            {users.map((user) => (
                <option key={user.id} value={user.login}>
                    [ {user.role} ] - {user.login}
                </option>
            ))}
        </select>
        <br></br>


            <label htmlFor="subsCount">count subs: {subsCount}</label>
            <input
                id="subsCount"
                type="range"
                min="0"
                max="50"
                value={subsCount}
                onChange={handleSubsChange}
            />


            <label htmlFor="daysCount">count days: {daysCount}</label>
            <input
                id="daysCount"
                type="range"
                min="0"
                max="30"
                value={daysCount}
                onChange={handleDaysChange}
            />

        <button onClick={handleSubmit}>gen subscribtion</button>

        <div>
            <h2>Keys list:</h2>
            {isLoading ? <p>Loading...</p> : (
                keys.map((key, index) => (
                    <div key={index}>
                        <p>{key.license_key}</p>
                    </div>
                ))
            )}
        </div>

        </div>
        )
    };



    const handlePrevPage = () => {
        if (currentPage > 1) setCurrentPage(currentPage - 1);
    };

    const handleNextPage = () => {
        setCurrentPage(currentPage + 1);
    };



    const handleFilterChange = (event) => {
        setFilter(event.target.value.toLowerCase());
    };

    // Функция для проверки вхождения фильтра в любое из полей объекта ключа лицензии
    const filterLicenseKeys = (key) => {
        return Object.values(key).some(value =>
            value.toString().toLowerCase().includes(filter)
        );
    };

    const filteredLicenseKeys = licenseKeys.filter(filterLicenseKeys);

    const renderInfoKeyList = () => {
      return (
          <div>
              <div>
                  <button onClick={handlePrevPage} disabled={currentPage === 1}>Предыдущая</button>
                  <button onClick={handleNextPage}>Следующая</button>
                   Страница {currentPage}
              </div>

              <input
                  type="text"
                  value={filter}
                  onChange={handleFilterChange}
                  placeholder="Поиск по всем полям..."
              />

              <table>
                  <thead>
                  <tr>
                      <th>ID</th>
                      <th>License Key</th>
                      <th>Cheat</th>
                      <th>TTL Cheat</th>
                      <th>Holder</th>
                      <th>Creator</th>
                      <th>Date of Creation</th>
                  </tr>
                  </thead>
                  <tbody>
                  {filteredLicenseKeys.map(key => (
                      <tr key={key.id}>
                          <td>{key.id}</td>
                          <td>{key.license_key}</td>
                          <td>{key.cheat}</td>
                          <td>{key.ttl_cheat}</td>
                          <td>{key.holder}</td>
                          <td>{key.creator}</td>
                          <td>{key.date_creation}</td>
                      </tr>
                  ))}
                  </tbody>
              </table>
          </div>
      )
    };

    return (
        <div>
            <h2>Cheat</h2>

            {renderGenerateCheat()}
            {renderInfoKeyList()}

        </div>
    );
}

export default LicenseKeys;