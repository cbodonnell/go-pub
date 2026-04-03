import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useEffect } from 'react'
import { AppLayout } from '@/components/layout/AppLayout'
import { HomePage } from '@/pages/HomePage'
import { ProfilePage } from '@/pages/ProfilePage'
import { PostPage } from '@/pages/PostPage'
import { ActivityPage } from '@/pages/ActivityPage'
import { NotFoundPage } from '@/pages/NotFoundPage'
import { useAuthStore } from '@/store/auth'
import { fetchCurrentUser } from '@/api/auth'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

function AuthInitializer({ children }: { children: React.ReactNode }) {
  const setUser = useAuthStore((s) => s.setUser)
  const setLoading = useAuthStore((s) => s.setLoading)

  useEffect(() => {
    fetchCurrentUser().then((user) => {
      setUser(user)
      setLoading(false)
    })
  }, [setUser, setLoading])

  return <>{children}</>
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthInitializer>
          <Routes>
            <Route element={<AppLayout />}>
              <Route index element={<HomePage />} />
              <Route path="/users/:name" element={<ProfilePage />} />
              <Route path="/actor" element={<ProfilePage />} />
              <Route path="/objects/:id" element={<PostPage />} />
              <Route path="/object" element={<PostPage />} />
              <Route path="/activities/:id" element={<ActivityPage />} />
              <Route path="/activity" element={<ActivityPage />} />
              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </Routes>
        </AuthInitializer>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
