module.exports = function (context, callback) {
    const size = context.request.body['size'] || 0;
    callback(200, new Buffer(size));
}

