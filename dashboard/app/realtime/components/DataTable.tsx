"use client"
import * as React from "react"
import {
    ColumnDef,
    flexRender,
    getCoreRowModel,
    ColumnFiltersState,
    SortingState,
    getSortedRowModel,
    getPaginationRowModel,
    VisibilityState,
    getFilteredRowModel,
    useReactTable,
} from "@tanstack/react-table"
import { Alert as AlertData } from "@/app/realtime/columns"
import { toast } from "sonner"

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
import { ArrowUpDown, MoreHorizontal } from "lucide-react"


import {
    Alert,
    AlertDescription,
    AlertTitle,
} from "@/components/ui/alert"
import {
    HoverCard,
    HoverCardContent,
    HoverCardTrigger,
} from "@/components/ui/hover-card"


import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { cloneUniformsGroups } from "three"

interface DataTableProps<TData, TValue> {
    columns: ColumnDef<TData, TValue>[]
    data: TData[]
    setData: React.Dispatch<React.SetStateAction<TData[]>>
    ackfn: (id: string) => Promise<boolean>
}

export function DataTable<TData, TValue>({
    columns,
    data,
    setData,
    ackfn,
}: DataTableProps<TData, TValue>) {
    const [sorting, setSorting] = React.useState<SortingState>([])
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
        []
    )
    const [rowSelection, setRowSelection] = React.useState({})
    const [columnVisibility, setColumnVisibility] =
        React.useState<VisibilityState>({})
    const table = useReactTable({
        data,
        columns,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        onColumnVisibilityChange: setColumnVisibility,
        getSortedRowModel: getSortedRowModel(),
        onRowSelectionChange: setRowSelection,
        state: { 
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
        },
    })

    const hancleOnClick = async () => {
        table.getSelectedRowModel().rows.forEach(async ({original}) => {
            console.log("Acknowledging alert")
            console.log(original)
            const id = (original as AlertData).id
            const acknowledged = await ackfn(id)
            if (acknowledged) {
                toast.success(`Alert ${id} acknowledged`)
                // data.forEach((alert) => {
                //     if ((alert as AlertData).id === id) {
                //         (alert as AlertData).Acknowledged = "1"
                //     }
                // })
                setData((prevData) => {
                    return prevData.map((alert) => {
                        if ((alert as AlertData).id === id) {
                            (alert as AlertData).Acknowledged = "1"
                        }
                        return alert
                    })
                })

            } else {
                toast.error(`Failed to acknowledge alert ${id}`)
            }
        })
        
        setRowSelection({})
    }

    return (
        <div>
            <div className="flex items-center py-4">
                <Input
                    placeholder="Filter alerts by category"
                    value={(table.getColumn("Category")?.getFilterValue() as string) ?? ""}
                    onChange={(event) =>
                        table.getColumn("Category")?.setFilterValue(event.target.value)
                    }
                    className="max-w-sm"
                />
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="outline" className="ml-auto">
                            Columns
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        {table
                            .getAllColumns()
                            .filter(
                                (column) => column.getCanHide()
                            )
                            .map((column) => {
                                return (
                                    <DropdownMenuCheckboxItem
                                        key={column.id}
                                        className="capitalize"
                                        checked={column.getIsVisible()}
                                        onCheckedChange={(value) =>
                                            column.toggleVisibility(!!value)
                                        }
                                    >
                                        {column.id}
                                    </DropdownMenuCheckboxItem>
                                )
                            })}
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        {table.getHeaderGroups().map((headerGroup) => (
                            <TableRow key={headerGroup.id}>
                                {headerGroup.headers.map((header) => {
                                    return (
                                        <TableHead key={header.id}>
                                            {header.isPlaceholder
                                                ? null
                                                : flexRender(
                                                    header.column.columnDef.header,
                                                    header.getContext()
                                                )}
                                        </TableHead>
                                    )
                                })}
                            </TableRow>
                        ))}
                    </TableHeader>
                    <TableBody>
                        {table.getRowModel().rows?.length ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow
                                    key={row.id}
                                    data-state={row.getIsSelected() && "selected"}
                                >
                                    {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id}>
                                            {flexRender(cell.column.columnDef.cell, {
                                                ...cell.getContext(),
                                                setData: setData,
                                            })}
                                        </TableCell>
                                        
                                    ))}
                                    {/* <TableCell>
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
                                                                                    Created At: {format(parse(row.getValue("CreatedAt"), "yyyy-MM-dd HH:mm:ss", new Date()), "yyyy-MM-dd HH:mm:ss")}
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
                                                        <AlertDialogAction>
                                                            <Button
                                                                onClick={async () => {
                                                                    console.log(`handling alert ${row.getValue("id")}`)
                                                                    const acknowledged = await hancleOnClick(row.getValue("id"))
                                                                    if (acknowledged) {
                                                                        toast.success(`Alert ${row.getValue("id")} acknowledged`)
                                                                    }
                                                                }}
                                                            >
                                                                Acknowledge
                                                            </Button>
                                                        </AlertDialogAction>
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
                                                                    Received At: {format(parse(row.getValue("CreatedAt"), "yyyy-MM-dd HH:mm:ss", new Date()), "yyyy-MM-dd HH:mm:ss")}
                                                                </span>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </HoverCardContent>
                                            </HoverCard>
                                        )}
                                    </TableCell> */}
                                    
                                </TableRow>
                            ))


                        ) : (
                            <TableRow>
                                <TableCell colSpan={columns.length} className="h-24 text-center">
                                    No results.
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            <div className="flex items-center justify-end space-x-2 py-4">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => table.previousPage()}
                    disabled={!table.getCanPreviousPage()}
                >
                    Previous
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => table.nextPage()}
                    disabled={!table.getCanNextPage()}
                >
                    Next
                </Button>
                <Button
                    variant="destructive"
                    size="sm"
                    onClick={async () => {
                        await hancleOnClick()
                    }}
                >
                    Acknowledge Selected
                </Button>
            </div>
        </div>
    )
}