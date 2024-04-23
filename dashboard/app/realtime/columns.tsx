import { ColumnDef } from "@tanstack/react-table"
import { z } from "zod"
import { cn } from "@/lib/utils";
import { parse, format, addDays, subMonths } from 'date-fns';
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
        header: "Category",
        accessorKey: "Category",
    },
    {
        header: "Severity",
        accessorKey: "Severity",
    },
    {
        header: "Source",
        accessorKey: "Source",
    },
    {
        header: "Created At",
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
            const acknowledged = row.getValue("Acknowledged")

            return (
                <span
                    className={cn(
                        "px-2 py-1 rounded-full text-xs font-semibold",
                        acknowledged === "0" ? "text-red-600" : "text-green-600"
                    )}
                >
                    {acknowledged === "0" ? "Pending" : "Acknowledged"}
                </span>
            )
        }
    },
    {
        header: "Remedy",
        accessorKey: "Remedy",
        cell: ({ row }) => {
            const dateString: string = row.getValue("CreatedAt")
            const date = parse(dateString, dateSchema, new Date())

            return (
                <div>
                    {row.getValue("Acknowledged") === "0" ? (
                        <AlertDialog>
                            <AlertDialogTrigger asChild>
                                <Button variant="outline">Acknowledge</Button>
                            </AlertDialogTrigger>
                            <AlertDialogContent>
                                <AlertDialogHeader>
                                    <AlertDialogTitle>Are you sure you want to acknowledge this alert?</AlertDialogTitle>
                                    <AlertDialogDescription className="flex flex-col gap-1 justify-center items-center">
                                        <Card className="w-full">
                                            <CardHeader>
                                                <CardTitle>{`Alert ${row.getValue("id")}`}</CardTitle>
                                            </CardHeader>
                                            <CardContent className="grid gap-4">
                                                <div className="grid grid-cols-2 gap-4">
                                                    <div className="space-y-1">
                                                        <p className="text-sm font-medium">Category</p>
                                                        <p className="text-gray-500 dark:text-gray-400">{row.getValue("Category")}</p>
                                                    </div>
                                                    <div className="space-y-1">
                                                        <p className="text-sm font-medium">Node</p>
                                                        <p className="text-gray-500 dark:text-gray-400">{row.getValue("node")}</p>
                                                    </div>
                                                </div>
                                                <div className="grid grid-cols-2 gap-4">
                                                    <div className="space-y-1">
                                                        <p className="text-sm font-medium">Severity</p>
                                                        <p className="text-gray-500 dark:text-gray-400">{row.getValue("Severity")}</p>
                                                    </div>
                                                    <div className="space-y-1">
                                                        <p className="text-sm font-medium">Source</p>
                                                        <p className="text-gray-500 dark:text-gray-400">{row.getValue("Source")}</p>
                                                    </div>
                                                </div>
                                            </CardContent>
                                            <CardFooter className="flex justify-end">
                                                <Alert variant="default">
                                                    <Activity className="h-4 w-4" />
                                                    <AlertTitle>Recommended Remedy!</AlertTitle>
                                                    <AlertDescription>
                                                        <div>
                                                            <p>
                                                                {row.getValue("Remedy")}
                                                            </p>
                                                            <p className="flex items-center pt-2">
                                                                <CalendarDays className="mr-2 h-6 w-6 opacity-70 rounded-lg" />{" "}
                                                                Created At: {format(date, "yyyy-MM-dd HH:mm:ss")}
                                                            </p>
                                                        </div>
                                                    </AlertDescription>
                                                </Alert>
                                            </CardFooter>
                                        </Card>
                                    </AlertDialogDescription>
                                </AlertDialogHeader>
                                <AlertDialogFooter>
                                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                                    <AlertDialogAction>Continue</AlertDialogAction>
                                </AlertDialogFooter>
                            </AlertDialogContent>
                        </AlertDialog>
                    ) : (
                        <HoverCard>
                            <HoverCardTrigger asChild>
                                <Button variant="outline" className="text-green-500">Acknowledged</Button>
                            </HoverCardTrigger>
                            <HoverCardContent className="w-80">
                                <div className="flex justify-between space-x-4">
                                    <Activity className="h-4 w-4" />
                                    <div className="space-y-1">
                                        <h4 className="text-sm font-semibold">Recommended Remedy</h4>
                                        <p className="text-sm">
                                            {row.getValue("Remedy")}
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
                    )}
                </div>
            )
        }
    }
]