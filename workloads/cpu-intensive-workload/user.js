module.exports = function(context, callback) {
    const ms = context.request.body['ms'] || 0;
    const start = Date.now();

    while(start + ms > Date.now());
    callback(200, `Calcuated for ${ms} ms`);
}

