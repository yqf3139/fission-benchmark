module.exports = function(context, callback) {
    const ms = context.request.body['size'] || 0;
    setTimeout(() => {
        callback(200, `Calcuated for ${ms} ms`);
    }, ms);
}

