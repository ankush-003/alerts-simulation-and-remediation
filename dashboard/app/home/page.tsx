"use client"
import React from 'react'
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from '@tanstack/react-query'
import CalendarStats from '@/components/CalendarStats'

const queryClient = new QueryClient()

function Home() {
  return (
    <QueryClientProvider client={queryClient}>
      <CalendarStats />
    </QueryClientProvider>
  )
}

export default Home