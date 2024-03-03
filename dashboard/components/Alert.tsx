import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { cn } from "@/lib/utils";

interface AlertProp {
    id: string;
    nodeID: string;
    description: string;
    severity: string;
    source: string;
    createdAt: string;
}

function Alert({ id, nodeID, description, severity, source, createdAt }: AlertProp) {
    return (
        <div>
            <Card>
                <CardHeader>
                    <CardTitle
                        className={cn(
                            "text-sm font-medium",
                            severity === "critical" && "text-red-500",
                            severity === "warning" && "text-white",
                            severity === "info" && "text-blue-500",
                        )}
                    >{`Alert ${id}`}</CardTitle>
                    <CardDescription>{`${description}`}</CardDescription>
                </CardHeader>
                <CardContent>
                    <p>{`NodeID : ${nodeID}`}</p><br />
                    <p>{`source : ${source}`}</p>
                </CardContent>
                <CardFooter>
                    <p>{`createdAt : ${createdAt}`}</p>
                </CardFooter>
            </Card>
        </div>
    )
}

export default Alert