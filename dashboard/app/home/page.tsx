"use client"
import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import CalendarStats from '@/components/charts/CalendarStats'
import SeverityStats from '@/components/charts/SeverityStats';

const queryClient = new QueryClient()

const Page = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className='w-full gap-2 grid'>
        <SeverityStats />
        <CalendarStats />
      </div>
    </QueryClientProvider>
  );
};

export default Page;