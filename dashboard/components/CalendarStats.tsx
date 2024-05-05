"use client"
import React from 'react'
// import { ResponsiveCalendar } from '@nivo/calendar'
import {
    useQuery,
} from '@tanstack/react-query'

import { BarChart, Bar, XAxis, YAxis, ResponsiveContainer, Tooltip, Label, CartesianGrid  } from 'recharts';

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
        <div className='w-full h-'>
            {/* <ResponsiveCalendar
        //         data={data}
        //         from="2024-01-01"
        //         to="2024-12-31"
        //         // emptyColor="#eeeeee"
        //         theme={
        //             {
        //                 background: `#0f0f0f`,
        //                 text: {
        //                     fill: '#ffffff',
        //                     fontSize: 11,
        //                 }
        //             }
        //         }
        //         // borderColor={{ theme: 'background' }}
        //         // colors={['#61cdbb', '#97e3d5', '#e8c1a0', '#f47560']}
        //         margin={{ top: 40, right: 40, bottom: 40, left: 40 }}
        //         yearSpacing={40}
        //         monthBorderColor="#ffffff"
        //         dayBorderWidth={2}
        //         dayBorderColor="#ffffff"
        //         legends={[
        //             {
        //                 anchor: 'bottom-right',
        //                 direction: 'row',
        //                 translateY: 36,
        //                 itemCount: 4,
        //                 itemWidth: 42,
        //                 itemHeight: 36,
        //                 itemsSpacing: 14,
        //                <ResponsiveLine
                data={
                    [
                        {
                            id: 'japan',
                            color: 'hsl(209, 70%, 50%)',
                            data: lineData
                        }
                    ]
                }
                margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                xScale={{ type: 'point' }}
                yScale={{
                    type: 'linear',
                    min: 'auto',
                    max: 'auto',
                    stacked: true,
                    reverse: false
                }}
                yFormat=" >-.2f"
                axisTop={null}
                axisRight={null}
                axisBottom={{
                    tickSize: 5,
                    tickPadding: 5,
                    tickRotation: 0,
                    legend: 'transportation',
                    legendOffset: 36,
                    legendPosition: 'middle',
                    truncateTickAt: 0
                }}
                axisLeft={{
                    tickSize: 5,
                    tickPadding: 5,
                    tickRotation: 0,
                    legend: 'count',
                    legendOffset: -40,
                    legendPosition: 'middle',
                    truncateTickAt: 0
                }}
                pointSize={10}
                pointColor={{ theme: 'background' }}
                pointBorderWidth={2}
                pointBorderColor={{ from: 'serieColor' }}
                pointLabel="data.yFormatted"
                pointLabelYOffset={-12}
                enableTouchCrosshair={true}
                useMesh={true}
                legends={[
                    {
                        anchor: 'bottom-right',
                        direction: 'column',
                        justify: false,
                        translateX: 100,
                        translateY: 0,
                        itemsSpacing: 0,
                        itemDirection: 'left-to-right',
                        itemWidth: 80,
                        itemHeight: 20,
                        itemOpacity: 0.75,
                        symbolSize: 12,
                        symbolShape: 'circle',
                        symbolBorderColor: 'rgba(0, 0, 0, .5)',
                        effects: [
                            {
                                on: 'hover',
                                style: {
                                    itemBackground: 'rgba(0, 0, 0, .03)',
                                    itemOpacity: 1
                                }
                            }
                        ]
                    }
                ]}
            />*/}

            <ResponsiveContainer width="100%" height={400}>
                <BarChart data={data} width={600} height={400}>
                {/* <CartesianGrid /> */}
                    <XAxis dataKey="day" >
                        <Label value="Your Alerts" offset={0} position="insideBottom" />
                        </XAxis>
                    <YAxis />
                    <Tooltip />
                    {/* <Legend /> */}
                    <Bar dataKey="value" barSize={30} fill="#ef4444">
                        {/* <LabelList dataKey="day" position="top" /> */}
                    </Bar>
                </BarChart>
            </ResponsiveContainer>
        </div>
    )
}