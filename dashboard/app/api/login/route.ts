import { NextResponse } from 'next/server';
import clientPromise from "../../../lib/mongodb";
import { compare } from 'bcrypt';

interface UserCredentials {
  email: string;
  password: string;
}

export async function POST(req: Request) {
  try {
    const { email, password } = await req.json() as UserCredentials;

    const client = await clientPromise;
    const db = client.db("users");
    const usersCollection = db.collection("details");

    const user = await usersCollection.findOne({ email });

    if (!user) {
      return NextResponse.json({ message: 'Invalid credentials' }, { status: 401 });
    }

    const isPasswordValid = await compare(password, user.password);

    if (!isPasswordValid) {
      return NextResponse.json({ message: 'Invalid credentials' }, { status: 401 });
    }

    return NextResponse.json({ message: 'Login successful' }, { status: 200 });
  } catch (e) {
    console.error(e);
    return NextResponse.json({ message: 'An error occurred' }, { status: 500 });
  }
}