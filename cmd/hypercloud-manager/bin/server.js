let express = require('express');
let config = require('config');

const options = {
    root: __dirname + '/../dist/',
    dotfiles: 'deny',
    headers: {}
};

var app = express();
app.disable('x-powered-by');

app.get('/health', (req, res) => {
    res.writeHead(200, {
        'Content-Type': 'application/json'
    });
    res.write(JSON.stringify({
        status: true
    }));
    res.end();
});

app.get('/*', (req, res, next) => {
    if (/\.(jpg|css|png|gif|ico|ttf|svg|woff2|eot|woff|js)$/i.test(req.path)) {
        res.sendFile(req.path, options, (err) => {
            if (err) {
                next(err);
            }
        });

        return;
    }

    res.sendFile('index.html', options, (err) => {
        if (err) {
            next(err);
        }
    });
});

app.listen(config.get('server.port'), () => {
    console.log(
        `Server is now running on http://localhost:${config.get('server.port')}/`
    );
});
