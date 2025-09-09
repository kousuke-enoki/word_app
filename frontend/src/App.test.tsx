import { render, screen } from '@testing-library/react'
import React from 'react'

import App from './App'

test('renders the main heading', () => {
  render(<App />)
  const linkElement = screen.getByText(/DictQuiz/i)
  expect(linkElement)
})
