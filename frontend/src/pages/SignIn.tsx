import React, { useState } from 'react';
import axios from 'axios';
// import { Link } from 'react-router-dom';

const SignIn: React.FC = () => {
  // const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

//   const handleSubmit = (event: React.FormEvent) => {
//     event.preventDefault();
//     // ログイン処理のロジックをここに追加
//     console.log('Email:', email);
//     console.log('Password:', password);
//   };

const handleSignIn = async (e: React.FormEvent<HTMLFormElement>) => {
  e.preventDefault();
  try {
    const response = await axios.post('http://localhost:8080/users/sign_in', {
        email,
        password,
    });
    setMessage('Sign in successful!');
    console.log(response)
  } catch (error) {
    setMessage('Sign in failed. Please try again.');
  }
};

  return (
    <div>
      <h1>サインイン</h1>
      <form onSubmit={handleSignIn}>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">サインイン</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
};

export default SignIn;