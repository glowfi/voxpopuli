import express from 'express';
import cors from 'cors';
import 'dotenv/config';
import { trpcExpress } from '@voxpopuli/trpc-server';

const app = express();
app.use(express.json());
app.use(cors());

app.get('/', (_, res) => {
    return res.json({ msg: 'Hello World!' }).status(200);
});

// TRPC Middleware
app.use('/trpc', trpcExpress);

app.listen(process.env.SERVER_PORT, () => {
    console.log(`Listening on http://localhost:${process.env.SERVER_PORT} ...`);
});
