import React from 'react'
import { Link } from 'react-router-dom'

const Header: React.FC = () => {
  return (
    <div>
      <h1>
        <Link to="/">word app</Link>
      </h1>
    </div>
  )
}

export default Header
