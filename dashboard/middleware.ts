import { NextResponse, NextRequest } from "next/server";
import { auth } from "@/auth";

export async function middleware(req: NextRequest) {
    const session = await auth();
    const url = req.nextUrl;
    
    if (session && (
        url.toString().startsWith("/")
    )) {
        console.log("redirecting to /home");
        return NextResponse.redirect("/home");
    }
    if (!session) 
        return NextResponse.redirect(new URL("/", req.url));

    return NextResponse.next();
}

export const config = {
    matcher: [
        '/home',
        '/alert-config',
        '/logs',
        '/about',
        '/nodes',
        '/realtime',
        // '/url',
        // '/url/:path*',
    ],
}