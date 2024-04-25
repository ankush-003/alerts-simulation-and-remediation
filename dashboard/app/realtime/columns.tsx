import { ColumnDef } from "@tanstack/react-table"
import { z } from "zod"
import { cn } from "@/lib/utils";
import { parse, format, addDays, subMonths, set } from 'date-fns';
import { Activity, CalendarDays } from 'lucide-react';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { CardTitle, CardHeader, CardContent, CardFooter, Card } from "@/components/ui/card"
import { ArrowUpDown, MoreHorizontal } from "lucide-react"
import { toast } from "sonner"
import { Checkbox } from "@/components/ui/checkbox"

import {
    Alert,
    AlertDescription,
    AlertTitle,
} from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import {
    HoverCard,
    HoverCardContent,
    HoverCardTrigger,
} from "@/components/ui/hover-card"
import { Console } from "console";
import { cloneUniformsGroups } from "three";
import { Badge } from "@/components/ui/badge"


const dateSchema = 'yyyy-MM-dd HH:mm:ss'

const alertSchema = z.object({
    id: z.string(),
    Category: z.string(),
    node: z.string(),
    Severity: z.string(),
    Source: z.string(),
    CreatedAt: z.string(),
    // received as "0" or "1" -> convert to boolean
    Acknowledged: z.string(),
    Remedy: z.string(),
})

const colorMap : Record<string, string> = {
    "Memory": "red",
    "CPU": "blue",
    "Disk": "green",
    "Network": "purple",
    "Power": "yellow",
    "Applications": "indigo",
    "Security": "pink",
    "RuntimeMetrics": "red",
    "Safe": "green",
    "Critical": "red",
    "Warning": "yellow",
}

export type Alert = z.infer<typeof alertSchema>;

export const columns: ColumnDef<Alert>[] = [
    // {  
    //     header: "Id",
    //     accessorKey: "id",
    // },  
    // {
    //     header: "Node Id",
    //     accessorKey: "node",
    // },
    {
        id: "select",
        header: ({ table }) => (
          <Checkbox
            checked={
              table.getIsAllPageRowsSelected() ||
              (table.getIsSomePageRowsSelected() && "indeterminate")
            }
            onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
            aria-label="Select all"
          />
        ),
        cell: ({ row }) => (
          <Checkbox
            checked={row.getIsSelected()}
            onCheckedChange={(value) => row.toggleSelected(!!value)}
            aria-label="Select row"
          />
        ),
    },
    {
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                >
                    Category
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
            )
        },
        accessorKey: "Category",
        cell: ({ row }) => {
            const alert = (row.original as Alert)
            return (
                <Badge
                    variant={
                        alert.Category === "RuntimeMetrics" ? "default" : "outline"
                    }
                >
                    {alert.Category}
                </Badge>
            )
        }

    },
    {
        header: "Severity",
        accessorKey: "Severity",
        cell: ({ row }) => {
            const alert = (row.original as Alert)
            return (
                <p className={`text-${colorMap[alert.Severity]}-500`}>
                    {alert.Severity}
                </p>
            )
        }
    },
    {
        header: "Source",
        accessorKey: "Source",
    },
    {
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                >
                    Created At
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
            )
        },
        accessorKey: "CreatedAt",
        cell: ({ row }) => {
            const dateString: string = row.getValue("CreatedAt")
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
        accessorKey: "Acknowledged",
        cell: ({ row }) => {
            const alert = (row.original as Alert)
            const dateString: string = alert.CreatedAt
            const date = parse(dateString, dateSchema, new Date())

            return (
                <div>
                    <HoverCard>
                            <HoverCardTrigger asChild>
                                <Button variant="outline" className={
                                    cn(
                                        "text-sm",
                                        row.getValue("Acknowledged") === '0' ? "text-red-600" : "text-green-600"
                                    )
                                }>
                                    {row.getValue("Acknowledged") === '0' ? "Pending" : "Acknowledged"}
                                </Button>
                            </HoverCardTrigger>
                            <HoverCardContent className="w-80">
                                <div className="flex justify-between space-x-4">
                                    <Activity className="h-4 w-4" />
                                    <div className="space-y-1">
                                        <h4 className="text-sm font-semibold">Recommended Remedy</h4>
                                        <p className="text-sm">
                                            {alert.Remedy}
                                        </p>
                                        <p>
                                            alert id: {alert.id}
                                        </p>
                                        <p>
                                            node id: {alert.node}
                                        </p>
                                        <div className="flex items-center pt-2">
                                            <CalendarDays className="mr-2 h-6 w-6 opacity-70 rounded-lg" />{" "}
                                            <span className="text-xs text-muted-foreground">
                                                Received At: {format(date, "yyyy-MM-dd HH:mm:ss")}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </HoverCardContent>
                        </HoverCard>
                </div>
            )
        }
    }
]