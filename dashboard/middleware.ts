import { NextResponse, NextRequest } from "next/server";
import { auth } from "@/auth";

export async function middleware(req: NextRequest) {
    const session = await auth();
    const url = req.nextUrl;
    
    if (session && (
        url.toString().startsWith("/")
    )) {
        return NextResponse.redirect(new URL("/home", req.url));
    }
    if (!session) 
        return NextResponse.redirect(new URL("/", req.url));

    return NextResponse.next();
}

export const config = {
    matcher: [
        '/home',
        '/profile',
        '/logs',
        // '/about',
        // '/nodes',
        // '/realtime',
        // '/chat',
        // '/url',
        // '/url/:path*',
    ],
}