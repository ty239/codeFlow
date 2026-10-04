import { useAuth } from './auth/context'
import { AuthPage } from './pages/AuthPage'
import { HomePage } from './pages/HomePage'

export default function App() {
  const { user } = useAuth()
  return user ? <HomePage /> : <AuthPage />
}
