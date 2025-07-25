import './App.css'
import './styles/global.css'

import React from 'react'

import AppRouter from './routes/AppRouter'

const App: React.FC = () => {
  return (
    <div className="App">
      <AppRouter />
    </div>
  )
}

export default App
