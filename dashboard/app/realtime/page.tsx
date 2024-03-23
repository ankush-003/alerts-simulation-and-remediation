"use client"

import { useEffect, useState } from "react"
import { Alert as AlertData, columns } from "./columns"
import { DataTable } from "@/components/DataTable";


export default function Realtime() {
    const [data, setData] = useState<AlertData[]>([]);

    useEffect(() => {
        console.log("Realtime page loaded")
        const eventSource = new EventSource("/api/stream")
        eventSource.addEventListener("message", (event) => {
            const data = JSON.parse(event.data)
            console.table(data)
            if(data.hasOwnProperty("message")) {
                console.log("Connection established")
                return
            }

            data.status = "open"
            setData((prevData) => [data, ...prevData])
        })

        return () => {
            eventSource.close()
        }
    }, []);

    return (
        <div className="container mx-auto py-10">
            <DataTable columns={columns} data={data} />
        </div>
    )
}