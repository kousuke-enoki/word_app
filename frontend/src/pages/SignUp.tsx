import React, { useState } from 'react';
import axios from 'axios';
// import { Link } from 'react-router-dom';

const SignUp: React.FC = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

//   const handleSubmit = (event: React.FormEvent) => {
//     event.preventDefault();
//     // ログイン処理のロジックをここに追加
//     console.log('Email:', email);
//     console.log('Password:', password);
//   };

const handleSignUp = async (e: React.FormEvent<HTMLFormElement>) => {
  e.preventDefault();
  try {
    const response = await axios.post('http://localhost:8080/users/sign_up', {
        name,
        email,
        password,
    });
    setMessage('Sign up successful!');
    console.log(response)
  } catch (error) {
    setMessage('Sign up failed. Please try again.');
  }
};

  return (
    <div>
      <h1>サインアップ</h1>
      <form onSubmit={handleSignUp}>
        <div>
          <label htmlFor="name">Name:</label>
          <input
            type="name"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
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
        <button type="submit">サインアップ</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
};

export default SignUp;