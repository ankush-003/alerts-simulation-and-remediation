"use client"
import React from 'react'
import { ResponsiveCalendar } from '@nivo/calendar'
import {
    useQuery,
  } from '@tanstack/react-query'

interface Data {
    day: string;
    value: number;
}

export default function CalendarStats() {
    const { isPending, error, data } = useQuery<Data[]>({
        queryKey: ['stats'],
        queryFn: async () => {
            const response = await fetch('/api/stats')
            return response.json()
        }
    })

    if (isPending) {
        return <div>Loading...</div>
    }

    if (error) {
        return <div>Error: {error.message}</div>
    }

    console.table(data)

    return (
        <div className='w-full h-[400px]'>
            <ResponsiveCalendar
                data={data}
                from="2024-01-01"
                to="2024-12-31"
                emptyColor="#eeeeee"
                theme={
                    {
                        background: '#000000',
                        text: {
                            fill: '#ffffff',
                            fontSize: 11,
                        }
                    }
                }
                // borderColor={{ theme: 'background' }}
                // colors={['#61cdbb', '#97e3d5', '#e8c1a0', '#f47560']}
                margin={{ top: 40, right: 40, bottom: 40, left: 40 }}
                yearSpacing={40}
                monthBorderColor="#ffffff"
                dayBorderWidth={2}
                dayBorderColor="#ffffff"
                legends={[
                    {
                        anchor: 'bottom-right',
                        direction: 'row',
                        translateY: 36,
                        itemCount: 4,
                        itemWidth: 42,
                        itemHeight: 36,
                        itemsSpacing: 14,
                        itemDirection: 'right-to-left'
                    }
                ]}
            />
        </div>
    )
}