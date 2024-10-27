import React from 'react'
import './App.css'
import AppRouter from './routes/AppRouter'
import './styles/global.css'

const App: React.FC = () => {
  return (
    <div className="App">
      <AppRouter />
    </div>
  )
}

export default App
