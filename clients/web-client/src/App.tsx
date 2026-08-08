import { useState } from 'react'
import { Navbar } from './components/Navbar'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
		Hello world
		<Navbar />
    </>
  )
}

export default App
