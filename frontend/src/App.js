import React, { useEffect, useState } from 'react';
import axios from 'axios';
import SignUp from './user/SignUp';
import SignIn from './user/SignIn';

function App() {
  const [data, setData] = useState(null);

  useEffect(() => {
    // APIリクエストを行う
    axios.get('http://localhost:8080/your-endpoint')
      .then(response => {
        setData(response.data);
      })
      .catch(error => {
        console.error('There was an error!', error);
      });
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Welcome to the App</h1>
        <div>
          <h2>Data from API</h2>
          {data ? <pre>{JSON.stringify(data, null, 2)}</pre> : <p>Loading...</p>}
        </div>
        <div>
          <SignUp />
          <SignIn />
        </div>
      </header>
    </div>
  );
}

export default App;
