import { ColumnDef } from "@tanstack/react-table"
import { z } from "zod"
import { cn } from "@/lib/utils";
import { parse, format, addDays, subMonths } from 'date-fns';

const dateSchema = 'yyyy-MM-dd HH:mm:ss'

const alertSchema = z.object({
    id: z.string(),
    node_id: z.string(),
    description: z.string(),
    severity: z.string(),
    source: z.string(),
    created_at: z.string(),
    status: z.enum(["open", "ack"]),
})

export type Alert = z.infer<typeof alertSchema>;

export const columns: ColumnDef<Alert>[] = [
    {
        header: "Id",
        accessorKey: "id",
    },
    {
        header: "Node Id",
        accessorKey: "node_id",
    },
    {
        header: "Description",
        accessorKey: "description",
    },
    {
        header: "Severity",
        accessorKey: "severity",
    },
    {
        header: "Source",
        accessorKey: "source",
    },
    {
        header: "Created At",
        accessorKey: "created_at",
        cell: ({ row }) => {
            const dateString:string = row.getValue("created_at")
            const date = parse(dateString, dateSchema, new Date())

            return (
                <span className="text-gray-500">
                    {format(date, "yyyy-MM-dd HH:mm:ss")}
                </span>
            )
        }
    },
    {
        header: "Status",
        accessorKey: "status",
    },
]