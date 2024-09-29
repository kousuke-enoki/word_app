import React, { useState } from 'react';
import axiosInstance from '../axiosConfig';

const SignUp: React.FC = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

const handleSignUp = async (e: React.FormEvent<HTMLFormElement>) => {
  e.preventDefault();
  try {
    const response = await axiosInstance.post('/users/sign_up', {
        name,
        email,
        password,
    });
    const token = response.data.token;
    localStorage.setItem('token', token);
    setMessage('Sign up successful!');
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