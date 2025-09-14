import express from 'express';
import type { Request, Response } from 'express';

const app = express();
const portArg = process.argv[2];
const PORT = portArg ? Number(portArg) : process.env.PORT || 3000;

app.get('/', (req: Request, res: Response) => {
    res.send("Request to index page");
});

app.get('/sample', (req: Request, res: Response) => {
    console.log("Received request to sample page");
});

app.get('/health', (req: Request, res: Response) => {
    res.send("OK");
});

// app.get('/[route]', (req: Request, res: Response) => {
//     res.send(`Request to dynamic route: ${req.params.route}`);
// });

app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});