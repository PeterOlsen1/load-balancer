import express from 'express';
var app = express();
var PORT = process.env.PORT || 3000;
app.get('/', function (req, res) {
    res.send("Request to index page");
});
app.get('/health', function (req, res) {
    res.send("OK");
});
// app.get('/[route]', (req: Request, res: Response) => {
//     res.send(`Request to dynamic route: ${req.params.route}`);
// });
app.listen(PORT, function () {
    console.log("Server is running on port ".concat(PORT));
});
