export const dynamic = 'force-dynamic'
import { Console } from 'console';
// export const runtime = 'edge';

// import { createClient } from 'redis';
import { Redis } from 'ioredis';

const redisClient = new Redis(process.env.REDIS_URL);

redisClient.on('error', (err) => console.error('Redis Client Error', err));
redisClient.on('connect', () => console.log('Connected to Redis'));

export async function GET() {
    const stream = 'alerts'; // Replace with your Redis Stream name
    const groupName = 'dashboard-group';
    let consumer;
    try {
        // try {
        //     await redisClient.xgroup('CREATE', stream, groupName, '$', 'MKSTREAM');
        // } catch (error) {
        //     console.error('Error creating group', error);
        // }
        const customReadable = new ReadableStream({
            async start(controller) {
                while (true) {
                    try {
                        consumer = await redisClient.xreadgroup('GROUP', groupName, 'dashboard-consumer', 'STREAMS', stream, '>');
                        if (consumer) {
                            console.log('Received messages', consumer);
                            const messages = consumer[0][1];
                            for (const message of messages) {
                                const [id, data] = message;
                                controller.enqueue(`data: ${JSON.stringify(data)}\n\n`);
                                await redisClient.xack(stream, groupName, id);
                            }
                        }
                    } catch (error) {
                        console.error('Error consuming messages', error);
                        break;
                    }
                }
            },
            async cancel() {
                // await consumer.disconnect();
                console.log('Stream closed');
            }
        });
        return new Response(customReadable, {
            headers: {
                'Content-Type': 'text/event-stream',
                Connection: 'keep-alive',
                'Cache-Control': 'no-cache, no-transform',
                'Content-Encoding': 'none'
            }
        });
    } catch (error) {
        console.error(error);
        // if (consumer) {
        //     await consumer.disconnect();
        // }
        return new Response("Internal Server Error", { status: 500 });
    }

}