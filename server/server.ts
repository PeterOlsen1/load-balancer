import express from 'express';
import type { Request, Response } from 'express';

const app = express();
const PORT = process.env.PORT || 3000;

app.get('/', (req: Request, res: Response) => {
    res.send("Request to index page");
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