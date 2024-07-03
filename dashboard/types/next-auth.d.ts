import 'next-auth'
import { DefaultSession } from 'next-auth'

declare module 'next-auth' {
    interface User {
        first_name: string,
        last_name: string,
        email: string,
        token: string,
        alert: any
    }
    interface Session {
        user: {
            first_name: string,
            last_name: string,
            email: string,
            token: string,
            alert: any
        } & DefaultSession['user']
    }
}

declare module 'next-auth/jwt' {
    interface JWT {
        first_name: string,
        last_name: string,
        email: string,
        token: string,
        alert: any
    }
}