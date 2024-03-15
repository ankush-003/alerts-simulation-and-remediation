"use client"

import { useEffect, useState } from "react"

interface Data {
    id: string
    node_id: string
    description: string
    severity: string
    source: string
    created_at: string    
}

export default function Realtime() {
    const [data, setData] = useState<Data[]>([]);

    useEffect(() => {
        console.log("Realtime page loaded")
        const eventSource = new EventSource("/api/stream")
        eventSource.addEventListener("message", (event) => {
            const data = JSON.parse(event.data)
            console.table(data)
            setData((prevData) => [data, ...prevData])
        })

        return () => {
            eventSource.close()
        }
    }, []);




    return (
        <div className="flex flex-col items-center justify-center h-[40rem]">
            <h1>Realtime</h1>
            <div>
                {/* create a dynamic table here */}
                <table>
                    <thead>
                        <tr>
                            <th>Id</th>
                            <th>Node Id</th>
                            <th>Description</th>
                            <th>Severity</th>
                            <th>Source</th>
                            <th>Created At</th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.map((item) => (
                            <tr key={item.id}>
                                <td>{item.id}</td>
                                <td>{item.node_id}</td>
                                <td>{item.description}</td>
                                <td>{item.severity}</td>
                                <td>{item.source}</td>
                                <td>{item.created_at}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    )
}