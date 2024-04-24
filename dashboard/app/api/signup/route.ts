import { NextResponse } from 'next/server';
import clientPromise from "../../../lib/mongodb";
import { hash } from 'argon2';

interface UserDetails {
  name: string;
  email: string;
  password: string;
}

export async function POST(req: Request) {
  try {
    const { name, email, password } = await req.json() as UserDetails;
    const hashedPassword = await hash(password);
    const client = await clientPromise;
    const db = client.db("AlertSimAndRemediation");
    const usersCollection = db.collection("Users");
    const result = await usersCollection.insertOne({
      name,
      email,
      password: hashedPassword,
    });
    if (result.acknowledged) {
      return NextResponse.json({ message: 'Signup successful' }, { status: 200 });
    } else {
      return NextResponse.json({ message: 'Signup failed' }, { status: 500 });
    }
  } catch (e) {
    console.error(e);
    return NextResponse.json({ message: 'An error occurred' }, { status: 500 });
  }
}