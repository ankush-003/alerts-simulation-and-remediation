// "use server";
import { MongoClient } from "mongodb";


// console.log("HEREE")


const uri = "mongodb+srv://alertssim:RRdxdKTHt99cEnWw@cluster0.aenmhq0.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";
const options = {};
const client =  new MongoClient(uri, options);
console.log("client", client)
let clientPromise: Promise<MongoClient>;


const getClient = async (): Promise<MongoClient> => {
  if (!client) {
    await client.connect();
  }
  return client;
};

// console.log("client",typeof(client))


clientPromise = getClient();

// console.log("ccccc", typeof(client))
// console.log(client._bsontype)



export default clientPromise;